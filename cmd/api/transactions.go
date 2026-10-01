package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

type transactionsResponse struct {
	AccountID    int64                `json:"account_id"`
	Transactions []domain.Transaction `json:"transactions"`
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

	default:
		writeJSON(
			w,
			http.StatusOK,
			transactionsResponse{
				AccountID:    accountID,
				Transactions: transactions,
			},
		)
	}
}
