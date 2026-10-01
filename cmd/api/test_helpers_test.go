package main

import (
	"context"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

type fakeAccountStore struct {
	createFunc func(
		ctx context.Context,
	) (*domain.Account, error)

	listFunc func(
		ctx context.Context,
	) ([]domain.Account, error)

	getByIDFunc func(
		ctx context.Context,
		id int64,
	) (*domain.Account, error)

	depositFunc func(
		ctx context.Context,
		accountID int64,
		amount int64,
	) (*domain.Account, error)

	transferFunc func(
		ctx context.Context,
		fromAccountID int64,
		toAccountID int64,
		amount int64,
	) (*domain.TransferResult, error)

	listTransactionsFunc func(
		ctx context.Context,
		accountID int64,
	) ([]domain.Transaction, error)
}

func (f *fakeAccountStore) Create(
	ctx context.Context,
) (*domain.Account, error) {
	return f.createFunc(ctx)
}

func (f *fakeAccountStore) List(
	ctx context.Context,
) ([]domain.Account, error) {
	return f.listFunc(ctx)
}

func (f *fakeAccountStore) GetByID(
	ctx context.Context,
	id int64,
) (*domain.Account, error) {
	return f.getByIDFunc(ctx, id)
}

func (f *fakeAccountStore) Deposit(
	ctx context.Context,
	accountID int64,
	amount int64,
) (*domain.Account, error) {
	return f.depositFunc(
		ctx,
		accountID,
		amount,
	)
}

func (f *fakeAccountStore) Transfer(
	ctx context.Context,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (*domain.TransferResult, error) {
	return f.transferFunc(
		ctx,
		fromAccountID,
		toAccountID,
		amount,
	)
}

func (f *fakeAccountStore) ListTransactions(
	ctx context.Context,
	accountID int64,
) ([]domain.Transaction, error) {
	return f.listTransactionsFunc(
		ctx,
		accountID,
	)
}
