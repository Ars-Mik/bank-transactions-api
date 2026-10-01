package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	app := &application{}

	request := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	response := httptest.NewRecorder()
	app.newRouter().ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusOK {

		t.Fatalf(
			"ожидался статус 200, получен %d",
			response.Code,
		)
	}

	expected := "{\"status\":\"ok\"}\n"

	if response.Body.String() != expected {

		t.Errorf(
			"ожидалось %q, получено %q",
			expected,
			response.Body.String(),
		)
	}
}
