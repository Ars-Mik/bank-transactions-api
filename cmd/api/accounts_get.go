package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func (app *application) getAccountHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil {

		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_ID",
			"Некорректный идентификатор счёта",
		)

		return
	}

	account, err := app.accounts.GetByID(
		r.Context(),
		id,
	)

	if err != nil {

		if err == sql.ErrNoRows {

			writeJSONError(
				w,
				http.StatusNotFound,
				"ACCOUNT_NOT_FOUND",
				"Банковский счёт не найден",
			)

			return
		}

		log.Printf(
			"Ошибка получения счёта: %v",
			err,
		)

		writeJSONError(
			w,
			http.StatusInternalServerError,
			"DATABASE_ERROR",
			"Не удалось получить счёт",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		account,
	)
}
