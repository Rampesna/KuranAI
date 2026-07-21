package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	appleVerifyProd = "https://buy.itunes.apple.com/verifyReceipt"
	appleVerifySand = "https://sandbox.itunes.apple.com/verifyReceipt"
)

// AppleReceiptService Apple receipt doğrulama.
type AppleReceiptService struct {
	sharedSecret string
	bundleID     string
}

// NewAppleReceiptService AppleReceiptService oluşturur.
func NewAppleReceiptService(sharedSecret, bundleID string) *AppleReceiptService {
	return &AppleReceiptService{sharedSecret: sharedSecret, bundleID: bundleID}
}

// VerifyRequest Apple verifyReceipt isteği.
type VerifyRequest struct {
	ReceiptData            string `json:"receipt-data"`
	Password               string `json:"password"`
	ExcludeOldTransactions bool   `json:"exclude-old-transactions,omitempty"`
}

// VerifyResponse Apple yanıtı.
type VerifyResponse struct {
	Status             int             `json:"status"`
	Receipt            *ReceiptInfo    `json:"receipt,omitempty"`
	LatestReceiptInfo  []InAppPurchase `json:"latest_receipt_info,omitempty"`
	PendingRenewalInfo []interface{}   `json:"pending_renewal_info,omitempty"`
}

type ReceiptInfo struct {
	BundleID string          `json:"bundle_id"`
	InApp    []InAppPurchase `json:"in_app,omitempty"`
}

type InAppPurchase struct {
	ProductID             string `json:"product_id"`
	TransactionID         string `json:"transaction_id"`
	OriginalTransactionID string `json:"original_transaction_id"`
	PurchaseDateMS        string `json:"purchase_date_ms"`
	ExpiresDateMS         string `json:"expires_date_ms"`
}

// VerifyResult doğrulama sonucu.
type VerifyResult struct {
	Valid         bool
	ProductID     string
	TransactionID string
	OriginalTxID  string
	ExpiresAtUnix int64
	IsSandbox     bool
}

// Verify receipt doğrular. Önce production dener, 21007 dönerse sandbox.
func (s *AppleReceiptService) Verify(receiptBase64 string) (*VerifyResult, error) {
	if s.sharedSecret == "" {
		return nil, fmt.Errorf("Apple IAP shared secret not configured")
	}
	reqBody := VerifyRequest{
		ReceiptData: receiptBase64,
		Password:    s.sharedSecret,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Production
	resp, err := http.Post(appleVerifyProd, "application/json", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var v VerifyResponse
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}

	isSandbox := false
	if v.Status == 21007 {
		isSandbox = true
		resp2, err := http.Post(appleVerifySand, "application/json", strings.NewReader(string(bodyBytes)))
		if err != nil {
			return nil, err
		}
		defer resp2.Body.Close()
		b, _ = io.ReadAll(resp2.Body)
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
	}

	if v.Status != 0 {
		return nil, fmt.Errorf("Apple verifyReceipt status: %d", v.Status)
	}
	if v.Receipt == nil || (s.bundleID != "" && v.Receipt.BundleID != s.bundleID) {
		return nil, fmt.Errorf("receipt bundle id is invalid")
	}

	inApps := v.LatestReceiptInfo
	if len(inApps) == 0 && v.Receipt != nil {
		inApps = v.Receipt.InApp
	}
	if len(inApps) == 0 {
		return nil, fmt.Errorf("no in-app purchases in receipt")
	}

	latest := inApps[0]
	var expiresMS int64
	for _, purchase := range inApps {
		candidate, _ := strconv.ParseInt(purchase.ExpiresDateMS, 10, 64)
		if candidate >= expiresMS {
			expiresMS = candidate
			latest = purchase
		}
	}
	txID := latest.OriginalTransactionID
	if txID == "" {
		txID = latest.TransactionID
	}

	return &VerifyResult{
		Valid:         true,
		ProductID:     latest.ProductID,
		TransactionID: latest.TransactionID,
		OriginalTxID:  txID,
		ExpiresAtUnix: expiresMS / 1000,
		IsSandbox:     isSandbox,
	}, nil
}
