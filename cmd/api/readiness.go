package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func (app *application) readinessHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		2*time.Second,
	)
	defer cancel()

	if err := app.db.PingContext(ctx); err != nil {
		log.Printf(
			"Проверка готовности PostgreSQL завершилась ошибкой: %v",
			err,
		)

		writeJSONError(
			w,
			http.StatusServiceUnavailable,
			"DATABASE_UNAVAILABLE",
			"Сервис временно не готов к работе",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status":   "ready",
			"database": "available",
		},
	)
}
