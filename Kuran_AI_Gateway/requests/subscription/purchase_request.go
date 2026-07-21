package subscription

// PurchaseRequest satın alma isteği (Apple receipt ile).
type PurchaseRequest struct {
	ReceiptData string `json:"receiptData" binding:"required"` // Base64 encoded receipt
}
