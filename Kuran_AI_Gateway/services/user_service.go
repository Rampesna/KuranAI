package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/rampesna/kuran-ai-gateway/core"
	"github.com/rampesna/kuran-ai-gateway/database"
	"github.com/rampesna/kuran-ai-gateway/models"
	userRequest "github.com/rampesna/kuran-ai-gateway/requests/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// usernameFromEmail e-posta adresinden kullanıcı adı üretir: @ öncesi kısımdan özel karakterleri kaldırır, _xxxxxx ekler.
func usernameFromEmail(email string) string {
	at := strings.Index(email, "@")
	local := "user"
	if at > 0 {
		local = email[:at]
	}
	var b strings.Builder
	for _, r := range local {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	base := strings.TrimSpace(b.String())
	if base == "" {
		base = "user"
	}
	return base + "_" + randomNumeric(6)
}

func randomNumeric(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%10]
	}
	return string(b)
}

func strPtr(s string) *string { return &s }

// UserService users tablosu için business logic.
type UserService struct {
	tokenService *TokenService
}

// NewUserService UserService oluşturur.
func NewUserService(tokenService *TokenService) *UserService {
	return &UserService{tokenService: tokenService}
}

// Me JWT'den gelen user ID ile mevcut kullanıcıyı döner.
func (s *UserService) Me(ctx context.Context, userID uint) *core.ServiceResponse {
	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(false, "User not found", 404, nil)
		}
		return core.NewServiceResponse(false, "Database error", 500, nil)
	}
	user.HasPassword = user.Password != ""
	return core.NewServiceResponse(true, "User fetched successfully", 200, user)
}

// Register kullanıcı kaydı.
func (s *UserService) Register(ctx context.Context, req *userRequest.RegisterRequest) *core.ServiceResponse {
	var exist models.User
	if err := database.DB.WithContext(ctx).Where("email = ?", req.Email).First(&exist).Error; err == nil {
		return core.NewServiceResponse(false, "Email already exists", 400, nil)
	}
	if err := database.DB.WithContext(ctx).Where("username = ?", req.Username).First(&exist).Error; err == nil {
		return core.NewServiceResponse(false, "Username already exists", 400, nil)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return core.NewServiceResponse(false, "Failed to hash password", 500, nil)
	}
	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashed),
		Name:     req.Name,
		Status:   1,
	}
	if err := database.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return core.NewServiceResponse(false, "Failed to create user", 500, nil)
	}
	user.HasPassword = user.Password != ""
	return core.NewServiceResponse(true, "User created successfully", 201, user)
}

// Login email/username + password ile giriş.
func (s *UserService) Login(ctx context.Context, req *userRequest.LoginRequest) *core.ServiceResponse {
	var user models.User
	if err := database.DB.WithContext(ctx).Where("email = ? OR username = ?", req.Credential, req.Credential).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(false, "Invalid credentials", 401, nil)
		}
		return core.NewServiceResponse(false, "Database error", 500, nil)
	}
	if user.Status != 1 {
		return core.NewServiceResponse(false, "Account is inactive", 403, nil)
	}
	if user.Password == "" {
		return core.NewServiceResponse(false, "Please use social login (Google/Apple)", 400, nil)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return core.NewServiceResponse(false, "Invalid credentials", 401, nil)
	}
	return s.tokenService.IssueTokens(ctx, user.ID, models.TokenableTypeUser, req.DeviceInfo)
}

// LoginWithGoogle Google ID token ile giriş/kayıt.
func (s *UserService) LoginWithGoogle(ctx context.Context, req *userRequest.GoogleLoginRequest) *core.ServiceResponse {
	payload, err := ValidateGoogleIDToken(ctx, req.IDToken)
	if err != nil {
		return core.NewServiceResponse(false, "Invalid Google token: "+err.Error(), 401, nil)
	}
	var user models.User
	if err := database.DB.WithContext(ctx).Where("google_id = ?", payload.Sub).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			email := payload.Email
			if email == "" {
				email = "google_" + payload.Sub + "@kuranai.local"
			}
			// E-posta zaten varsa (önceki kayıt/deneme) mevcut hesaba google_id bağla
			if err := database.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err == nil {
				if err := database.DB.WithContext(ctx).Model(&user).Update("google_id", strPtr(payload.Sub)).Error; err != nil {
					return core.NewServiceResponse(false, "Failed to link Google account", 500, nil)
				}
			} else {
				username := usernameFromEmail(email)
				name := strings.TrimSpace(payload.Name)
				if name == "" {
					name = username
				}
				user = models.User{
					Email:    email,
					GoogleID: strPtr(payload.Sub),
					AppleID:  nil,
					Name:     name,
					Username: username,
					Status:   1,
				}
				if err := database.DB.WithContext(ctx).Create(&user).Error; err != nil {
					return core.NewServiceResponse(false, "Failed to create user", 500, nil)
				}
			}
		} else {
			return core.NewServiceResponse(false, "Database error", 500, nil)
		}
	}
	if user.Status != 1 {
		return core.NewServiceResponse(false, "Account is inactive", 403, nil)
	}
	return s.tokenService.IssueTokens(ctx, user.ID, models.TokenableTypeUser, req.DeviceInfo)
}

// LoginWithApple Apple identity token ile giriş/kayıt.
func (s *UserService) LoginWithApple(ctx context.Context, req *userRequest.AppleLoginRequest) *core.ServiceResponse {
	payload, err := ValidateAppleIDToken(ctx, req.IdentityToken)
	if err != nil {
		return core.NewServiceResponse(false, "Invalid Apple token: "+err.Error(), 401, nil)
	}
	var user models.User
	if err := database.DB.WithContext(ctx).Where("apple_id = ?", payload.Sub).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			email := payload.Email
			if email == "" {
				email = "apple_" + payload.Sub + "@kuranai.local"
			}
			// E-posta zaten varsa mevcut hesaba apple_id bağla
			if err := database.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err == nil {
				if err := database.DB.WithContext(ctx).Model(&user).Update("apple_id", strPtr(payload.Sub)).Error; err != nil {
					return core.NewServiceResponse(false, "Failed to link Apple account", 500, nil)
				}
			} else {
				username := usernameFromEmail(email)
				name := strings.TrimSpace(payload.Email)
				if name == "" {
					name = username
				}
				user = models.User{
					Email:    email,
					GoogleID: nil,
					AppleID:  strPtr(payload.Sub),
					Name:     name,
					Username: username,
					Status:   1,
				}
				if err := database.DB.WithContext(ctx).Create(&user).Error; err != nil {
					return core.NewServiceResponse(false, "Failed to create user", 500, nil)
				}
			}
		} else {
			return core.NewServiceResponse(false, "Database error", 500, nil)
		}
	}
	if user.Status != 1 {
		return core.NewServiceResponse(false, "Account is inactive", 403, nil)
	}
	return s.tokenService.IssueTokens(ctx, user.ID, models.TokenableTypeUser, req.DeviceInfo)
}

// UpdateUsername kullanıcı adını günceller; başkası kullanıyorsa hata döner.
func (s *UserService) UpdateUsername(ctx context.Context, userID uint, username string) *core.ServiceResponse {
	username = strings.TrimSpace(username)
	if len(username) < 2 || len(username) > 50 {
		return core.NewServiceResponse(false, "Username must be 2-50 characters", 400, nil)
	}
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return core.NewServiceResponse(false, "Username can only contain letters, numbers and underscore", 400, nil)
		}
	}
	var exist models.User
	if err := database.DB.WithContext(ctx).Where("username = ? AND id != ?", username, userID).First(&exist).Error; err == nil {
		return core.NewServiceResponse(false, "Username already taken", 400, nil)
	}
	if err := database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("username", username).Error; err != nil {
		return core.NewServiceResponse(false, "Failed to update username", 500, nil)
	}
	var user models.User
	_ = database.DB.WithContext(ctx).First(&user, userID)
	user.HasPassword = user.Password != ""
	return core.NewServiceResponse(true, "Username updated", 200, user)
}

// UploadAvatar profil resmini kaydeder; dosya yolunu döner.
func (s *UserService) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) *core.ServiceResponse {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		return core.NewServiceResponse(false, "Only jpg, png, gif, webp allowed", 400, nil)
	}
	if file.Size > 5*1024*1024 {
		return core.NewServiceResponse(false, "File too large (max 5MB)", 400, nil)
	}
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return core.NewServiceResponse(false, "Failed to create upload dir", 500, nil)
	}
	filename := fmt.Sprintf("%d%s", userID, ext)
	dst := filepath.Join(uploadDir, filename)
	_ = os.Remove(dst)
	src, err := file.Open()
	if err != nil {
		return core.NewServiceResponse(false, "Failed to read file", 500, nil)
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return core.NewServiceResponse(false, "Failed to save file", 500, nil)
	}
	defer out.Close()
	if _, err := out.ReadFrom(src); err != nil {
		os.Remove(dst)
		return core.NewServiceResponse(false, "Failed to write file", 500, nil)
	}
	avatarURL := "/uploads/avatars/" + filename
	if err := database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error; err != nil {
		os.Remove(dst)
		return core.NewServiceResponse(false, "Failed to update user", 500, nil)
	}
	var user models.User
	_ = database.DB.WithContext(ctx).First(&user, userID)
	user.HasPassword = user.Password != ""
	return core.NewServiceResponse(true, "Avatar uploaded", 200, map[string]interface{}{"avatarUrl": avatarURL, "user": user})
}

// RemoveAvatar profil resmini siler ve avatar_url'i temizler.
func (s *UserService) RemoveAvatar(ctx context.Context, userID uint) *core.ServiceResponse {
	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(false, "User not found", 404, nil)
		}
		return core.NewServiceResponse(false, "Database error", 500, nil)
	}
	if user.AvatarURL != "" {
		_ = os.Remove("." + user.AvatarURL)
	}
	if err := database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("avatar_url", "").Error; err != nil {
		return core.NewServiceResponse(false, "Failed to remove avatar", 500, nil)
	}
	user.AvatarURL = ""
	user.HasPassword = user.Password != ""
	return core.NewServiceResponse(true, "Avatar removed", 200, user)
}

// UpdatePassword mevcut şifreyi doğrular ve yeni şifreyle günceller.
// Social login (Google/Apple) ile kayıt olmuş kullanıcıların şifresi boş olabilir;
// bu durumda currentPassword kontrolü atlanır ve doğrudan yeni şifre set edilir.
func (s *UserService) UpdatePassword(ctx context.Context, userID uint, currentPassword, newPassword string) *core.ServiceResponse {
	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(false, "User not found", 404, nil)
		}
		return core.NewServiceResponse(false, "Database error", 500, nil)
	}

	// Kullanıcının mevcut şifresi varsa eski şifreyi doğrula
	if user.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
			return core.NewServiceResponse(false, "Current password is incorrect", 401, nil)
		}
	}
	// Şifresi yoksa (social login) — eski şifre sorulmaz, doğrudan yeni şifre set edilir

	if len(newPassword) < 6 {
		return core.NewServiceResponse(false, "Password must be at least 6 characters", 400, nil)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return core.NewServiceResponse(false, "Failed to hash password", 500, nil)
	}
	if err := database.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", string(hashed)).Error; err != nil {
		return core.NewServiceResponse(false, "Failed to update password", 500, nil)
	}
	return core.NewServiceResponse(true, "Password updated", 200, nil)
}

// DeleteAccount kullanıcının hesabını ve ilişkili verilerini kalıcı olarak siler (Apple Guideline 5.1.1(v)).
func (s *UserService) DeleteAccount(ctx context.Context, userID uint) *core.ServiceResponse {
	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.NewServiceResponse(false, "User not found", 404, nil)
		}
		return core.NewServiceResponse(false, "Database error", 500, nil)
	}

	// İlişkili verileri sil (soft/hard delete)
	tx := database.DB.WithContext(ctx).Begin()

	// Token'ları sil
	tx.Where("tokenable_id = ?", userID).Delete(&models.JWTToken{})
	// Sohbet mesajlarını sil
	tx.Exec("DELETE FROM chat_messages WHERE session_id IN (SELECT id FROM chat_sessions WHERE user_id = ?)", userID)
	// Sohbet oturumlarını sil
	tx.Where("user_id = ?", userID).Delete(&models.ChatSession{})
	// Kullanım loglarını sil
	tx.Where("user_id = ?", userID).Delete(&models.UsageLog{})
	// Abonelikleri sil
	tx.Where("user_id = ?", userID).Delete(&models.UserSubscription{})
	// Avatar dosyasını sil
	if user.AvatarURL != "" {
		_ = os.Remove("." + user.AvatarURL)
	}
	// Kullanıcıyı kalıcı olarak sil (hard delete)
	if err := tx.Unscoped().Delete(&user).Error; err != nil {
		tx.Rollback()
		return core.NewServiceResponse(false, "Failed to delete account", 500, nil)
	}

	tx.Commit()
	return core.NewServiceResponse(true, "Account deleted permanently", 200, nil)
}
