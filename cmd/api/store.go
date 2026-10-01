package main

import (
	"context"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

type accountStore interface {
	Create(
		ctx context.Context,
	) (*domain.Account, error)

	List(
		ctx context.Context,
	) ([]domain.Account, error)

	GetByID(
		ctx context.Context,
		id int64,
	) (*domain.Account, error)

	Deposit(
		ctx context.Context,
		accountID int64,
		amount int64,
	) (*domain.Account, error)

	Transfer(
		ctx context.Context,
		fromAccountID int64,
		toAccountID int64,
		amount int64,
	) (*domain.TransferResult, error)

	ListTransactions(
		ctx context.Context,
		accountID int64,
	) ([]domain.Transaction, error)
}
