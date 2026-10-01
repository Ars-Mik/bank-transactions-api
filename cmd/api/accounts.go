package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

type accountsResponse struct {
	Accounts []domain.Account `json:"accounts"`
}

func (app *application) createAccountHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	account, err := app.accounts.Create(r.Context())

	if err != nil {

		log.Printf(
			"Ошибка создания счёта: %v",
			err,
		)

		writeJSONError(
			w,
			http.StatusInternalServerError,
			"DATABASE_ERROR",
			"Не удалось создать банковский счёт",
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(account); err != nil {

		log.Printf(
			"Ошибка отправки ответа: %v",
			err,
		)
	}
}

func (app *application) listAccountsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	accounts, err := app.accounts.List(r.Context())
	if err != nil {
		log.Printf("Ошибка получения списка счетов: %v", err)

		writeJSONError(
			w,
			http.StatusInternalServerError,
			"DATABASE_ERROR",
			"Не удалось получить список банковских счетов",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		accountsResponse{
			Accounts: accounts,
		},
	)
}
