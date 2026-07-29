package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rampesna/kuran-ai-gateway/config"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

var (
	ErrGooglePlayUnavailable = errors.New("google play verification unavailable")
	ErrInvalidGooglePurchase = errors.New("invalid google play purchase")
	ErrGoogleAccountConflict = errors.New("google play purchase belongs to another account")
)

type GooglePlayPurchaseInfo struct {
	PackageName         string
	ProductID           string
	PurchaseToken       string
	OrderID             string
	ObfuscatedAccountID string
	Environment         string
	PurchasedAt         time.Time
	ExpiresAt           time.Time
	Acknowledged        bool
	RawJSON             []byte
}

type GooglePlayPurchaseVerifier interface {
	VerifySubscription(ctx context.Context, packageName, productID, purchaseToken string) (GooglePlayPurchaseInfo, error)
	Ready() bool
}

type googlePlayAPIVerifier struct {
	service     *androidpublisher.Service
	packageName string
	allowTest   bool
}

func NewGooglePlayVerifier(cfg config.GooglePlayConfig) (GooglePlayPurchaseVerifier, error) {
	credentialsPath := strings.TrimSpace(cfg.CredentialsPath)
	if credentialsPath == "" {
		return nil, nil
	}
	if _, err := os.Stat(credentialsPath); err != nil {
		return nil, fmt.Errorf("google play credentials: %w", err)
	}
	service, err := androidpublisher.NewService(
		context.Background(),
		option.WithCredentialsFile(credentialsPath),
	)
	if err != nil {
		return nil, fmt.Errorf("google play client: %w", err)
	}
	packageName := strings.TrimSpace(cfg.PackageName)
	if packageName == "" {
		packageName = "com.talhacan.kuranai"
	}
	return &googlePlayAPIVerifier{
		service:     service,
		packageName: packageName,
		allowTest:   cfg.AllowTest,
	}, nil
}

func (v *googlePlayAPIVerifier) Ready() bool {
	return v != nil && v.service != nil
}

func (v *googlePlayAPIVerifier) VerifySubscription(
	ctx context.Context,
	packageName string,
	productID string,
	purchaseToken string,
) (GooglePlayPurchaseInfo, error) {
	if !v.Ready() {
		return GooglePlayPurchaseInfo{}, ErrGooglePlayUnavailable
	}

	packageName = strings.TrimSpace(packageName)
	if packageName == "" {
		packageName = v.packageName
	}
	productID = strings.TrimSpace(productID)
	purchaseToken = strings.TrimSpace(purchaseToken)
	if packageName != v.packageName || productID == "" || purchaseToken == "" {
		return GooglePlayPurchaseInfo{}, ErrInvalidGooglePurchase
	}

	subscription, err := v.service.Purchases.Subscriptionsv2.
		Get(packageName, purchaseToken).
		Context(ctx).
		Do()
	if err != nil {
		return GooglePlayPurchaseInfo{}, fmt.Errorf("%w: %v", ErrInvalidGooglePurchase, err)
	}
	rawJSON, _ := json.Marshal(subscription)

	info := GooglePlayPurchaseInfo{
		PackageName:   packageName,
		PurchaseToken: purchaseToken,
		OrderID:       strings.TrimSpace(subscription.LatestOrderId),
		Environment:   "Production",
		Acknowledged: strings.EqualFold(
			subscription.AcknowledgementState,
			"ACKNOWLEDGEMENT_STATE_ACKNOWLEDGED",
		),
		RawJSON: rawJSON,
	}
	if subscription.ExternalAccountIdentifiers != nil {
		info.ObfuscatedAccountID = strings.TrimSpace(
			subscription.ExternalAccountIdentifiers.ObfuscatedExternalAccountId,
		)
	}
	if subscription.TestPurchase != nil {
		info.Environment = "Sandbox"
	}
	if info.Environment == "Sandbox" && !v.allowTest {
		return GooglePlayPurchaseInfo{}, fmt.Errorf("%w: test purchase rejected", ErrInvalidGooglePurchase)
	}
	if purchasedAt, parseErr := time.Parse(time.RFC3339, strings.TrimSpace(subscription.StartTime)); parseErr == nil {
		info.PurchasedAt = purchasedAt.UTC()
	}

	for _, item := range subscription.LineItems {
		if item == nil || strings.TrimSpace(item.ProductId) != productID {
			continue
		}
		info.ProductID = productID
		if expiresAt, parseErr := time.Parse(time.RFC3339, strings.TrimSpace(item.ExpiryTime)); parseErr == nil {
			info.ExpiresAt = expiresAt.UTC()
		}
		break
	}

	if info.ProductID == "" || info.ExpiresAt.IsZero() {
		return GooglePlayPurchaseInfo{}, fmt.Errorf("%w: product or expiry mismatch", ErrInvalidGooglePurchase)
	}
	if info.PurchasedAt.IsZero() {
		info.PurchasedAt = time.Now().UTC()
	}
	switch strings.ToUpper(strings.TrimSpace(subscription.SubscriptionState)) {
	case "SUBSCRIPTION_STATE_ACTIVE",
		"SUBSCRIPTION_STATE_IN_GRACE_PERIOD",
		"SUBSCRIPTION_STATE_CANCELED":
		// Canceled subscriptions remain entitled until their verified expiry.
	default:
		if !info.ExpiresAt.After(time.Now().UTC()) {
			return GooglePlayPurchaseInfo{}, fmt.Errorf("%w: subscription inactive", ErrInvalidGooglePurchase)
		}
	}
	if !info.ExpiresAt.After(time.Now().UTC()) {
		return GooglePlayPurchaseInfo{}, fmt.Errorf("%w: subscription expired", ErrInvalidGooglePurchase)
	}
	return info, nil
}
