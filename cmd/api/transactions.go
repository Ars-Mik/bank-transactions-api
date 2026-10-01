package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

const (
	transactionDirectionIncoming = "incoming"
	transactionDirectionOutgoing = "outgoing"
)

type transactionHistoryItem struct {
	ID            int64     `json:"id"`
	Kind          string    `json:"kind"`
	Direction     string    `json:"direction"`
	AmountKopecks int64     `json:"amount_kopecks"`
	FromAccountID *int64    `json:"from_account_id"`
	ToAccountID   *int64    `json:"to_account_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type transactionsResponse struct {
	AccountID    int64                    `json:"account_id"`
	Transactions []transactionHistoryItem `json:"transactions"`
}

func transactionToHistoryItem(
	transaction domain.Transaction,
	accountID int64,
) (transactionHistoryItem, error) {

	direction, err := transactionDirection(
		transaction,
		accountID,
	)

	if err != nil {
		return transactionHistoryItem{}, err
	}

	return transactionHistoryItem{
		ID:            transaction.ID,
		Kind:          transaction.Kind,
		Direction:     direction,
		AmountKopecks: transaction.AmountKopecks,
		FromAccountID: transaction.FromAccountID,
		ToAccountID:   transaction.ToAccountID,
		CreatedAt:     transaction.CreatedAt,
	}, nil
}

func transactionDirection(
	transaction domain.Transaction,
	accountID int64,
) (string, error) {

	switch transaction.Kind {

	case "deposit":

		if transaction.ToAccountID != nil &&
			*transaction.ToAccountID == accountID {

			return transactionDirectionIncoming, nil
		}

	case "transfer":

		if transaction.FromAccountID != nil &&
			*transaction.FromAccountID == accountID {

			return transactionDirectionOutgoing, nil
		}

		if transaction.ToAccountID != nil &&
			*transaction.ToAccountID == accountID {

			return transactionDirectionIncoming, nil
		}
	}

	return "", fmt.Errorf(
		"транзакция %d не относится к счёту %d",
		transaction.ID,
		accountID,
	)
}

func (app *application) accountTransactionsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	accountID, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || accountID <= 0 {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_ACCOUNT_ID",
			"Некорректный идентификатор банковского счёта",
		)
		return
	}

	transactions, err := app.accounts.ListTransactions(
		r.Context(),
		accountID,
	)

	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		writeJSONError(
			w,
			http.StatusNotFound,
			"ACCOUNT_NOT_FOUND",
			"Банковский счёт не найден",
		)
		return

	case err != nil:
		log.Printf(
			"Ошибка получения истории счёта %d: %v",
			accountID,
			err,
		)

		writeJSONError(
			w,
			http.StatusInternalServerError,
			"DATABASE_ERROR",
			"Не удалось получить историю операций",
		)
		return
	}

	items := make(
		[]transactionHistoryItem,
		0,
		len(transactions),
	)

	for _, transaction := range transactions {

		item, err := transactionToHistoryItem(
			transaction,
			accountID,
		)

		if err != nil {
			log.Printf(
				"Ошибка формирования истории счёта %d: %v",
				accountID,
				err,
			)

			writeJSONError(
				w,
				http.StatusInternalServerError,
				"TRANSACTION_DATA_ERROR",
				"Не удалось сформировать историю операций",
			)
			return
		}

		items = append(items, item)
	}

	writeJSON(
		w,
		http.StatusOK,
		transactionsResponse{
			AccountID:    accountID,
			Transactions: items,
		},
	)
}
