package domain

import "time"

type Transaction struct {
	ID            int64     `json:"id"`
	Kind          string    `json:"kind"`
	AmountKopecks int64     `json:"amount_kopecks"`
	FromAccountID *int64    `json:"from_account_id"`
	ToAccountID   *int64    `json:"to_account_id"`
	CreatedAt     time.Time `json:"created_at"`
}
