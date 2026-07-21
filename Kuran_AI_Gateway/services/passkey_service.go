package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/rampesna/kuran-ai-gateway/config"
	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/dto"
	"github.com/rampesna/kuran-ai-gateway/models"
	redisstore "github.com/rampesna/kuran-ai-gateway/redis"
	"gorm.io/gorm"
)

const passkeyCeremonyTTL = 5 * time.Minute

type passkeyCeremony struct {
	Kind    string               `json:"kind"`
	UserID  uint                 `json:"userId,omitempty"`
	Session webauthn.SessionData `json:"session"`
}

type passkeyUser struct {
	model       models.User
	credentials []webauthn.Credential
}

func (u *passkeyUser) WebAuthnID() []byte                         { return u.model.PasskeyUserHandle }
func (u *passkeyUser) WebAuthnName() string                       { return u.model.Email }
func (u *passkeyUser) WebAuthnDisplayName() string                { return u.model.Name }
func (u *passkeyUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

// PasskeyService implements WebAuthn passkey registration and discoverable login.
// Ceremony state lives in Redis so requests may land on either gateway pod.
type PasskeyService struct {
	webAuthn     *webauthn.WebAuthn
	tokenService *TokenService
}

func NewPasskeyService(appConfig *config.Config, tokenService *TokenService) (*PasskeyService, error) {
	w, err := webauthn.New(&webauthn.Config{
		RPID:          appConfig.WebAuthn.RPID,
		RPDisplayName: appConfig.WebAuthn.RPName,
		RPOrigins:     appConfig.WebAuthn.RPOrigins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			ResidentKey:        protocol.ResidentKeyRequirementRequired,
			RequireResidentKey: protocol.ResidentKeyRequired(),
			UserVerification:   protocol.VerificationRequired,
		},
		Timeouts: webauthn.TimeoutsConfig{
			Login:        webauthn.TimeoutConfig{Enforce: true, Timeout: passkeyCeremonyTTL},
			Registration: webauthn.TimeoutConfig{Enforce: true, Timeout: passkeyCeremonyTTL},
		},
	})
	if err != nil {
		return nil, err
	}
	return &PasskeyService{webAuthn: w, tokenService: tokenService}, nil
}

func (s *PasskeyService) BeginRegistration(ctx context.Context, userID uint) *core.ServiceResponse {
	user, err := s.loadUserByID(ctx, userID, true)
	if err != nil {
		return passkeyServiceError(err)
	}
	options, session, err := s.webAuthn.BeginRegistration(
		user,
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
		webauthn.WithExclusions(webauthn.Credentials(user.credentials).CredentialDescriptors()),
	)
	if err != nil {
		return core.NewServiceResponse(false, "Could not create passkey challenge", 500, nil)
	}
	transactionID, err := s.saveCeremony(ctx, passkeyCeremony{Kind: "register", UserID: userID, Session: *session})
	if err != nil {
		return core.NewServiceResponse(false, "Could not save passkey challenge", 500, nil)
	}
	return core.NewServiceResponse(true, "Passkey registration challenge created", 200, map[string]interface{}{
		"transactionId": transactionID,
		"options":       options.Response,
	})
}

func (s *PasskeyService) FinishRegistration(ctx context.Context, userID uint, transactionID string, request *http.Request) *core.ServiceResponse {
	ceremony, err := s.consumeCeremony(ctx, transactionID)
	if err != nil || ceremony.Kind != "register" || ceremony.UserID != userID {
		return core.NewServiceResponse(false, "Passkey challenge is invalid or expired", 400, nil)
	}
	user, err := s.loadUserByID(ctx, userID, false)
	if err != nil {
		return passkeyServiceError(err)
	}
	credential, err := s.webAuthn.FinishRegistration(user, ceremony.Session, request)
	if err != nil {
		return core.NewServiceResponse(false, "Passkey registration could not be verified", 400, nil)
	}
	encoded, err := json.Marshal(credential)
	if err != nil {
		return core.NewServiceResponse(false, "Could not save passkey", 500, nil)
	}
	record := models.PasskeyCredential{
		UserID:         userID,
		CredentialID:   credential.ID,
		CredentialData: encoded,
	}
	if err := database.DB.WithContext(ctx).Create(&record).Error; err != nil {
		return core.NewServiceResponse(false, "Passkey is already registered or could not be saved", 409, nil)
	}
	return core.NewServiceResponse(true, "Passkey registered", 201, map[string]interface{}{"count": len(user.credentials) + 1})
}

func (s *PasskeyService) BeginLogin(ctx context.Context) *core.ServiceResponse {
	options, session, err := s.webAuthn.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationRequired))
	if err != nil {
		return core.NewServiceResponse(false, "Could not create passkey challenge", 500, nil)
	}
	transactionID, err := s.saveCeremony(ctx, passkeyCeremony{Kind: "login", Session: *session})
	if err != nil {
		return core.NewServiceResponse(false, "Could not save passkey challenge", 500, nil)
	}
	return core.NewServiceResponse(true, "Passkey login challenge created", 200, map[string]interface{}{
		"transactionId": transactionID,
		"options":       options.Response,
	})
}

func (s *PasskeyService) FinishLogin(ctx context.Context, transactionID string, request *http.Request, device *dto.DeviceInfo) *core.ServiceResponse {
	ceremony, err := s.consumeCeremony(ctx, transactionID)
	if err != nil || ceremony.Kind != "login" {
		return core.NewServiceResponse(false, "Passkey challenge is invalid or expired", 400, nil)
	}
	resolvedUser, credential, err := s.webAuthn.FinishPasskeyLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		user, loadErr := s.loadUserByHandle(ctx, userHandle)
		if loadErr != nil {
			return nil, loadErr
		}
		for _, owned := range user.credentials {
			if string(owned.ID) == string(rawID) {
				return user, nil
			}
		}
		return nil, errors.New("credential does not belong to user")
	}, ceremony.Session, request)
	if err != nil {
		return core.NewServiceResponse(false, "Passkey login could not be verified", 401, nil)
	}
	user, ok := resolvedUser.(*passkeyUser)
	if !ok || user.model.Status != 1 {
		return core.NewServiceResponse(false, "Account is inactive or unavailable", 403, nil)
	}
	encoded, err := json.Marshal(credential)
	if err != nil {
		return core.NewServiceResponse(false, "Could not update passkey", 500, nil)
	}
	now := time.Now()
	if err := database.DB.WithContext(ctx).Model(&models.PasskeyCredential{}).
		Where("user_id = ? AND credential_id = ?", user.model.ID, credential.ID).
		Updates(map[string]interface{}{"credential_data": encoded, "last_used_at": &now}).Error; err != nil {
		return core.NewServiceResponse(false, "Could not update passkey", 500, nil)
	}
	return s.tokenService.IssueTokens(ctx, user.model.ID, models.TokenableTypeUser, device)
}

func (s *PasskeyService) Status(ctx context.Context, userID uint) *core.ServiceResponse {
	var count int64
	if err := database.DB.WithContext(ctx).Model(&models.PasskeyCredential{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return core.NewServiceResponse(false, "Could not read passkey status", 500, nil)
	}
	return core.NewServiceResponse(true, "Passkey status fetched", 200, map[string]interface{}{"count": count, "enabled": count > 0})
}

func (s *PasskeyService) DeleteAll(ctx context.Context, userID uint) *core.ServiceResponse {
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.PasskeyCredential{}).Error; err != nil {
		return core.NewServiceResponse(false, "Could not remove passkeys", 500, nil)
	}
	return core.NewServiceResponse(true, "Passkeys removed", 200, map[string]interface{}{"count": 0, "enabled": false})
}

func (s *PasskeyService) loadUserByID(ctx context.Context, userID uint, ensureHandle bool) (*passkeyUser, error) {
	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	if ensureHandle && len(user.PasskeyUserHandle) == 0 {
		handle := make([]byte, 32)
		if _, err := rand.Read(handle); err != nil {
			return nil, err
		}
		if err := database.DB.WithContext(ctx).Model(&user).Update("passkey_user_handle", handle).Error; err != nil {
			return nil, err
		}
		user.PasskeyUserHandle = handle
	}
	credentials, err := loadPasskeyCredentials(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &passkeyUser{model: user, credentials: credentials}, nil
}

func (s *PasskeyService) loadUserByHandle(ctx context.Context, handle []byte) (*passkeyUser, error) {
	var user models.User
	if len(handle) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err := database.DB.WithContext(ctx).Where("passkey_user_handle = ?", handle).First(&user).Error; err != nil {
		return nil, err
	}
	credentials, err := loadPasskeyCredentials(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &passkeyUser{model: user, credentials: credentials}, nil
}

func loadPasskeyCredentials(ctx context.Context, userID uint) ([]webauthn.Credential, error) {
	var records []models.PasskeyCredential
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&records).Error; err != nil {
		return nil, err
	}
	credentials := make([]webauthn.Credential, 0, len(records))
	for _, record := range records {
		var credential webauthn.Credential
		if err := json.Unmarshal(record.CredentialData, &credential); err != nil {
			return nil, fmt.Errorf("decode credential: %w", err)
		}
		credentials = append(credentials, credential)
	}
	return credentials, nil
}

func (s *PasskeyService) saveCeremony(ctx context.Context, ceremony passkeyCeremony) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	transactionID := base64.RawURLEncoding.EncodeToString(randomBytes)
	encoded, err := json.Marshal(ceremony)
	if err != nil {
		return "", err
	}
	if err := redisstore.Client.Set(ctx, redisstore.KeyPrefixPasskey+transactionID, encoded, passkeyCeremonyTTL).Err(); err != nil {
		return "", err
	}
	return transactionID, nil
}

func (s *PasskeyService) consumeCeremony(ctx context.Context, transactionID string) (*passkeyCeremony, error) {
	if len(transactionID) < 32 || len(transactionID) > 128 {
		return nil, errors.New("invalid transaction")
	}
	encoded, err := redisstore.Client.GetDel(ctx, redisstore.KeyPrefixPasskey+transactionID).Bytes()
	if err != nil {
		return nil, err
	}
	var ceremony passkeyCeremony
	if err := json.Unmarshal(encoded, &ceremony); err != nil {
		return nil, err
	}
	return &ceremony, nil
}

func passkeyServiceError(err error) *core.ServiceResponse {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.NewServiceResponse(false, "User not found", 404, nil)
	}
	return core.NewServiceResponse(false, "Passkey service error", 500, nil)
}
