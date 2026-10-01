package domain

import (
	"errors"
	"math"
)

var (
	ErrInvalidAmount     = errors.New("сумма должна быть больше нуля")
	ErrBalanceOverflow   = errors.New("превышен максимально допустимый баланс")
	ErrAccountNotFound   = errors.New("счёт не найден")
	ErrInsufficientFunds = errors.New("недостаточно средств на счёте")
	ErrSameAccount       = errors.New("нельзя выполнить перевод на тот же счёт")
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

func (a *Account) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if a.BalanceKopecks < amount {
		return ErrInsufficientFunds
	}

	a.BalanceKopecks -= amount

	return nil
}
