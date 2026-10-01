//go:build integration

package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestAccountRepositoryTransferIntegration(
	t *testing.T,
) {
	dsn := os.Getenv("TEST_DATABASE_URL")

	if dsn == "" {
		t.Skip(
			"TEST_DATABASE_URL не задан: интеграционный тест пропущен",
		)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf(
			"не удалось открыть подключение к тестовой БД: %v",
			err,
		)
	}
	defer db.Close()

	pingContext, cancelPing := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	if err := db.PingContext(pingContext); err != nil {
		cancelPing()

		t.Fatalf(
			"тестовая PostgreSQL недоступна: %v",
			err,
		)
	}

	cancelPing()

	repository := NewAccountRepository(db)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	fromAccount, err := repository.Create(ctx)
	if err != nil {
		t.Fatalf(
			"не удалось создать счёт отправителя: %v",
			err,
		)
	}

	toAccount, err := repository.Create(ctx)
	if err != nil {
		t.Fatalf(
			"не удалось создать счёт получателя: %v",
			err,
		)
	}

	t.Cleanup(func() {
		cleanupContext, cancelCleanup := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancelCleanup()

		_, _ = db.ExecContext(
			cleanupContext,
			`
				DELETE FROM transactions
				WHERE from_account_id IN ($1, $2)
				   OR to_account_id IN ($1, $2)
			`,
			fromAccount.ID,
			toAccount.ID,
		)

		_, _ = db.ExecContext(
			cleanupContext,
			`
				DELETE FROM accounts
				WHERE id IN ($1, $2)
			`,
			fromAccount.ID,
			toAccount.ID,
		)
	})

	if _, err := repository.Deposit(
		ctx,
		fromAccount.ID,
		100000,
	); err != nil {
		t.Fatalf(
			"не удалось пополнить счёт отправителя: %v",
			err,
		)
	}

	if _, err := repository.Deposit(
		ctx,
		toAccount.ID,
		50000,
	); err != nil {
		t.Fatalf(
			"не удалось пополнить счёт получателя: %v",
			err,
		)
	}

	beforeFrom, err := repository.GetByID(
		ctx,
		fromAccount.ID,
	)
	if err != nil {
		t.Fatalf(
			"не удалось получить счёт отправителя: %v",
			err,
		)
	}

	beforeTo, err := repository.GetByID(
		ctx,
		toAccount.ID,
	)
	if err != nil {
		t.Fatalf(
			"не удалось получить счёт получателя: %v",
			err,
		)
	}

	totalBefore :=
		beforeFrom.BalanceKopecks +
			beforeTo.BalanceKopecks

	result, err := repository.Transfer(
		ctx,
		fromAccount.ID,
		toAccount.ID,
		30000,
	)
	if err != nil {
		t.Fatalf(
			"перевод завершился ошибкой: %v",
			err,
		)
	}

	if result.TransactionID <= 0 {
		t.Errorf(
			"ожидался положительный ID транзакции, получен %d",
			result.TransactionID,
		)
	}

	afterFrom, err := repository.GetByID(
		ctx,
		fromAccount.ID,
	)
	if err != nil {
		t.Fatalf(
			"не удалось получить баланс отправителя после перевода: %v",
			err,
		)
	}

	afterTo, err := repository.GetByID(
		ctx,
		toAccount.ID,
	)
	if err != nil {
		t.Fatalf(
			"не удалось получить баланс получателя после перевода: %v",
			err,
		)
	}

	if afterFrom.BalanceKopecks != 70000 {
		t.Errorf(
			"ожидался баланс отправителя 70000, получен %d",
			afterFrom.BalanceKopecks,
		)
	}

	if afterTo.BalanceKopecks != 80000 {
		t.Errorf(
			"ожидался баланс получателя 80000, получен %d",
			afterTo.BalanceKopecks,
		)
	}

	totalAfter :=
		afterFrom.BalanceKopecks +
			afterTo.BalanceKopecks

	if totalAfter != totalBefore {
		t.Errorf(
			"нарушен денежный инвариант: до перевода %d, после %d",
			totalBefore,
			totalAfter,
		)
	}

	t.Run(
		"failed transfer does not change balances",
		func(t *testing.T) {
			errResult, err := repository.Transfer(
				ctx,
				fromAccount.ID,
				toAccount.ID,
				1000000,
			)

			if errResult != nil {
				t.Errorf(
					"при ошибочном переводе результат должен быть nil",
				)
			}

			if !errors.Is(
				err,
				domain.ErrInsufficientFunds,
			) {
				t.Fatalf(
					"ожидалась ErrInsufficientFunds, получена %v",
					err,
				)
			}

			currentFrom, err := repository.GetByID(
				ctx,
				fromAccount.ID,
			)
			if err != nil {
				t.Fatal(err)
			}

			currentTo, err := repository.GetByID(
				ctx,
				toAccount.ID,
			)
			if err != nil {
				t.Fatal(err)
			}

			if currentFrom.BalanceKopecks !=
				afterFrom.BalanceKopecks {

				t.Errorf(
					"баланс отправителя изменился после неудачного перевода",
				)
			}

			if currentTo.BalanceKopecks !=
				afterTo.BalanceKopecks {

				t.Errorf(
					"баланс получателя изменился после неудачного перевода",
				)
			}
		},
	)
}
