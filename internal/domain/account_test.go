package domain

import (
	"errors"
	"math"
	"testing"
)

func TestDeposit(t *testing.T) {
	tests := []struct {
		name           string
		initialBalance int64
		amount         int64
		wantBalance    int64
		wantErr        error
	}{
		{
			name:           "successful deposit",
			initialBalance: 0,
			amount:         150000,
			wantBalance:    150000,
		},
		{
			name:           "zero amount",
			initialBalance: 100000,
			amount:         0,
			wantBalance:    100000,
			wantErr:        ErrInvalidAmount,
		},
		{
			name:           "negative amount",
			initialBalance: 100000,
			amount:         -100,
			wantBalance:    100000,
			wantErr:        ErrInvalidAmount,
		},
		{
			name:           "balance overflow",
			initialBalance: math.MaxInt64,
			amount:         1,
			wantBalance:    math.MaxInt64,
			wantErr:        ErrBalanceOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := Account{
				ID:             1,
				BalanceKopecks: tt.initialBalance,
			}

			err := account.Deposit(tt.amount)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"expected error %v, got %v",
					tt.wantErr,
					err,
				)
			}

			if account.BalanceKopecks != tt.wantBalance {
				t.Errorf(
					"expected balance %d, got %d",
					tt.wantBalance,
					account.BalanceKopecks,
				)
			}
		})
	}
}
