package subscription

type GooglePurchaseRequest struct {
	PackageName   string `json:"packageName"`
	ProductID     string `json:"productId" binding:"required"`
	PurchaseToken string `json:"purchaseToken" binding:"required"`
}
