package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

func TestTransferInsufficientFunds(
	t *testing.T,
) {
	store := &fakeAccountStore{
		transferFunc: func(
			ctx context.Context,
			fromAccountID int64,
			toAccountID int64,
			amount int64,
		) (*domain.TransferResult, error) {
			return nil, domain.ErrInsufficientFunds
		},
	}

	app := &application{
		accounts: store,
	}

	body := `
	{
		"from_account_id": 1,
		"to_account_id": 2,
		"amount_kopecks": 500000
	}
	`

	request := httptest.NewRequest(
		http.MethodPost,
		"/transfers",
		strings.NewReader(body),
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	response := httptest.NewRecorder()

	app.newRouter().ServeHTTP(
		response,
		request,
	)

	if response.Code !=
		http.StatusUnprocessableEntity {

		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusUnprocessableEntity,
			response.Code,
		)
	}

	if !strings.Contains(
		response.Body.String(),
		`"code":"INSUFFICIENT_FUNDS"`,
	) {
		t.Errorf(
			"получен неожиданный ответ: %s",
			response.Body.String(),
		)
	}
}
