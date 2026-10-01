package domain

import (
	"errors"
	"math"
)

var (
	ErrInvalidAmount   = errors.New("Сумма должна быть больше нуля")
	ErrBalanceOverflow = errors.New("Превышен максимально допустимый баланс")
	ErrAccountNotFound = errors.New("Счёт не найден")
)

type Account struct {
	ID             int64 `json:"id"`
	BalanceKopecks int64 `json:"balance_kopecks"`
}

func (a *Account) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if amount > math.MaxInt64-a.BalanceKopecks {
		return ErrBalanceOverflow
	}

	a.BalanceKopecks += amount

	return nil
}
