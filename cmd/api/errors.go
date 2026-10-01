package main

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error struct {
		Code string `json:"code"`

		Message string `json:"message"`
	} `json:"error"`
}

func writeJSONError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {

	response := errorResponse{}

	response.Error.Code = code
	response.Error.Message = message

	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)
}
