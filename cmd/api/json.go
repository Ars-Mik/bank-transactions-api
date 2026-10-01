package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)

		http.Error(
			w,
			`{"error":{"code":"JSON_ENCODING_ERROR","message":"Не удалось сформировать ответ"}}`,
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(status)

	if _, err := w.Write(append(payload, '\n')); err != nil {
		log.Printf("Ошибка отправки HTTP-ответа: %v", err)
	}
}
