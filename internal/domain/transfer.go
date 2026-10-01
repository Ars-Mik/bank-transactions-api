package domain

type TransferResult struct {
	TransactionID int64   `json:"transaction_id"`
	FromAccount   Account `json:"from_account"`
	ToAccount     Account `json:"to_account"`
}
