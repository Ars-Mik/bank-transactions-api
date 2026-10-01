package main

import (
	"testing"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

func int64Pointer(value int64) *int64 {
	return &value
}

func TestTransactionDirection(t *testing.T) {
	tests := []struct {
		name        string
		transaction domain.Transaction
		accountID   int64
		want        string
		wantErr     bool
	}{
		{
			name: "deposit is incoming",
			transaction: domain.Transaction{
				ID:          1,
				Kind:        "deposit",
				ToAccountID: int64Pointer(1),
			},
			accountID: 1,
			want:      transactionDirectionIncoming,
		},
		{
			name: "transfer is outgoing for sender",
			transaction: domain.Transaction{
				ID:            2,
				Kind:          "transfer",
				FromAccountID: int64Pointer(1),
				ToAccountID:   int64Pointer(2),
			},
			accountID: 1,
			want:      transactionDirectionOutgoing,
		},
		{
			name: "transfer is incoming for recipient",
			transaction: domain.Transaction{
				ID:            2,
				Kind:          "transfer",
				FromAccountID: int64Pointer(1),
				ToAccountID:   int64Pointer(2),
			},
			accountID: 2,
			want:      transactionDirectionIncoming,
		},
		{
			name: "unrelated account returns error",
			transaction: domain.Transaction{
				ID:            2,
				Kind:          "transfer",
				FromAccountID: int64Pointer(1),
				ToAccountID:   int64Pointer(2),
			},
			accountID: 3,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := transactionDirection(
				tt.transaction,
				tt.accountID,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatal(
						"ожидалась ошибка, но она не была получена",
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"неожиданная ошибка: %v",
					err,
				)
			}

			if got != tt.want {
				t.Errorf(
					"ожидалось направление %q, получено %q",
					tt.want,
					got,
				)
			}
		})
	}
}
