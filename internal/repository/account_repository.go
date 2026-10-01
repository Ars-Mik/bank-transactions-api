package repository

import (
	"context"
	"database/sql"
	"errors"

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

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}

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

func (r *AccountRepository) ListTransactions(
	ctx context.Context,
	accountID int64,
) ([]domain.Transaction, error) {
	const existsQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM accounts
			WHERE id = $1
		)
	`

	var exists bool

	if err := r.db.QueryRowContext(
		ctx,
		existsQuery,
		accountID,
	).Scan(&exists); err != nil {
		return nil, err
	}

	if !exists {
		return nil, domain.ErrAccountNotFound
	}

	const query = `
		SELECT
			id,
			kind,
			amount_kopecks,
			from_account_id,
			to_account_id,
			created_at
		FROM transactions
		WHERE from_account_id = $1
		   OR to_account_id = $1
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]domain.Transaction, 0)

	for rows.Next() {
		var transaction domain.Transaction

		var fromAccountID sql.NullInt64
		var toAccountID sql.NullInt64

		if err := rows.Scan(
			&transaction.ID,
			&transaction.Kind,
			&transaction.AmountKopecks,
			&fromAccountID,
			&toAccountID,
			&transaction.CreatedAt,
		); err != nil {
			return nil, err
		}

		if fromAccountID.Valid {
			value := fromAccountID.Int64
			transaction.FromAccountID = &value
		}

		if toAccountID.Valid {
			value := toAccountID.Int64
			transaction.ToAccountID = &value
		}

		transactions = append(
			transactions,
			transaction,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *AccountRepository) Deposit(
	ctx context.Context,
	accountID int64,
	amount int64,
) (*domain.Account, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	account := &domain.Account{}

	const selectQuery = `
		SELECT id, balance_kopecks
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	err = tx.QueryRowContext(
		ctx,
		selectQuery,
		accountID,
	).Scan(
		&account.ID,
		&account.BalanceKopecks,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	if err := account.Deposit(amount); err != nil {
		return nil, err
	}

	const updateQuery = `
		UPDATE accounts
		SET balance_kopecks = $1
		WHERE id = $2
	`

	if _, err := tx.ExecContext(
		ctx,
		updateQuery,
		account.BalanceKopecks,
		account.ID,
	); err != nil {
		return nil, err
	}

	const transactionQuery = `
		INSERT INTO transactions (
			kind,
			amount_kopecks,
			from_account_id,
			to_account_id
		)
		VALUES (
			'deposit',
			$1,
			NULL,
			$2
		)
	`

	if _, err := tx.ExecContext(
		ctx,
		transactionQuery,
		amount,
		account.ID,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) Transfer(
	ctx context.Context,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (*domain.TransferResult, error) {
	if fromAccountID == toAccountID {
		return nil, domain.ErrSameAccount
	}

	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	firstID := fromAccountID
	secondID := toAccountID

	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	firstAccount, err := getAccountForUpdate(
		ctx,
		tx,
		firstID,
	)
	if err != nil {
		return nil, err
	}

	secondAccount, err := getAccountForUpdate(
		ctx,
		tx,
		secondID,
	)
	if err != nil {
		return nil, err
	}

	var fromAccount *domain.Account
	var toAccount *domain.Account

	if firstAccount.ID == fromAccountID {
		fromAccount = firstAccount
		toAccount = secondAccount
	} else {
		fromAccount = secondAccount
		toAccount = firstAccount
	}

	if err := fromAccount.Withdraw(amount); err != nil {
		return nil, err
	}

	if err := toAccount.Deposit(amount); err != nil {
		return nil, err
	}

	const updateQuery = `
		UPDATE accounts
		SET balance_kopecks = $1
		WHERE id = $2
	`

	if _, err := tx.ExecContext(
		ctx,
		updateQuery,
		fromAccount.BalanceKopecks,
		fromAccount.ID,
	); err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(
		ctx,
		updateQuery,
		toAccount.BalanceKopecks,
		toAccount.ID,
	); err != nil {
		return nil, err
	}

	const transactionQuery = `
		INSERT INTO transactions (
			kind,
			amount_kopecks,
			from_account_id,
			to_account_id
		)
		VALUES (
			'transfer',
			$1,
			$2,
			$3
		)
		RETURNING id
	`

	var transactionID int64

	err = tx.QueryRowContext(
		ctx,
		transactionQuery,
		amount,
		fromAccount.ID,
		toAccount.ID,
	).Scan(&transactionID)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.TransferResult{
		TransactionID: transactionID,
		FromAccount:   *fromAccount,
		ToAccount:     *toAccount,
	}, nil
}

func getAccountForUpdate(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
) (*domain.Account, error) {
	const query = `
		SELECT id, balance_kopecks
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	account := &domain.Account{}

	err := tx.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&account.ID,
		&account.BalanceKopecks,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return account, nil
}
