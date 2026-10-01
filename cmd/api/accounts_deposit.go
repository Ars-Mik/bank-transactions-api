package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

type depositRequest struct {
	AmountKopecks int64 `json:"amount_kopecks"`
}

func (app *application) depositHandler(
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

	var input depositRequest

	if err := readJSON(w, r, &input); err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_JSON",
			"Некорректные данные запроса",
		)

		return
	}

	account, err := app.accounts.Deposit(
		r.Context(),
		accountID,
		input.AmountKopecks,
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

	case errors.Is(err, domain.ErrInvalidAmount):
		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_AMOUNT",
			"Сумма пополнения должна быть больше нуля",
		)

		return

	case errors.Is(err, domain.ErrBalanceOverflow):
		writeJSONError(
			w,
			http.StatusUnprocessableEntity,
			"BALANCE_OVERFLOW",
			"Сумма пополнения превышает допустимый баланс",
		)

		return

	case err != nil:
		log.Printf(
			"Ошибка пополнения счёта %d: %v",
			accountID,
			err,
		)

		writeJSONError(
			w,
			http.StatusInternalServerError,
			"DATABASE_ERROR",
			"Не удалось пополнить банковский счёт",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		account,
	)
}
