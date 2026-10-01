package repository

import (
	"context"
	"database/sql"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Create(ctx context.Context) (*domain.Account, error) {
	const query = `
		INSERT INTO accounts (balance_kopecks)
		VALUES (0)
		RETURNING id, balance_kopecks
	`

	account := &domain.Account{}

	err := r.db.QueryRowContext(ctx, query).Scan(
		&account.ID,
		&account.BalanceKopecks,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) GetByID(
	ctx context.Context,
	id int64,
) (*domain.Account, error) {

	const query = `SELECT id, balance_kopecks FROM accounts
		 WHERE id = $1`

	account := &domain.Account{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&account.ID,
		&account.BalanceKopecks,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) List(
	ctx context.Context,
) ([]domain.Account, error) {
	const query = `
		SELECT id, balance_kopecks
		FROM accounts
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]domain.Account, 0)

	for rows.Next() {
		var account domain.Account

		if err := rows.Scan(
			&account.ID,
			&account.BalanceKopecks,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
