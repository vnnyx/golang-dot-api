package web

type TransactionCreateRequest struct {
	Name   string `json:"name"`
	UserID string
}

type TransactionUpdateRequest struct {
	TransactionID string
	Name          string `json:"name"`
}

type TransactionResponse struct {
	TransactionID string `json:"transaction_id"`
	Name          string `json:"name"`
	UserID        string `json:"user_id"`
}
