package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const maxRequestBodySize = 1 << 20 // 1 MiB

func readJSON(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) error {
	r.Body = http.MaxBytesReader(
		w,
		r.Body,
		maxRequestBodySize,
	)

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("тело запроса должно содержать один JSON-объект")
	}

	return nil
}

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
