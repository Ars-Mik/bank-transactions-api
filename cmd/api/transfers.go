package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	AmountKopecks int64 `json:"amount_kopecks"`
}

func (app *application) transferHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input transferRequest

	if err := readJSON(w, r, &input); err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_JSON",
			"Некорректные данные запроса",
		)
		return
	}

	if input.FromAccountID <= 0 || input.ToAccountID <= 0 {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_ACCOUNT_ID",
			"Некорректный идентификатор банковского счёта",
		)
		return
	}

	result, err := app.accounts.Transfer(
		r.Context(),
		input.FromAccountID,
		input.ToAccountID,
		input.AmountKopecks,
	)

	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		writeJSONError(
			w,
			http.StatusNotFound,
			"ACCOUNT_NOT_FOUND",
			"Один из банковских счетов не найден",
		)

	case errors.Is(err, domain.ErrSameAccount):
		writeJSONError(
			w,
			http.StatusBadRequest,
			"SAME_ACCOUNT",
			"Нельзя переводить средства на тот же счёт",
		)

	case errors.Is(err, domain.ErrInvalidAmount):
		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_AMOUNT",
			"Сумма перевода должна быть больше нуля",
		)

	case errors.Is(err, domain.ErrInsufficientFunds):
		writeJSONError(
			w,
			http.StatusUnprocessableEntity,
			"INSUFFICIENT_FUNDS",
			"Недостаточно средств на счёте",
		)

	case errors.Is(err, domain.ErrBalanceOverflow):
		writeJSONError(
			w,
			http.StatusUnprocessableEntity,
			"BALANCE_OVERFLOW",
			"Баланс получателя превышает допустимое значение",
		)

	case err != nil:
		log.Printf("Ошибка перевода: %v", err)

		writeJSONError(
			w,
			http.StatusInternalServerError,
			"DATABASE_ERROR",
			"Не удалось выполнить перевод",
		)

	default:
		writeJSON(
			w,
			http.StatusOK,
			result,
		)
	}
}
