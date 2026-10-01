package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ars-Mik/bank-transactions-api/internal/domain"
)

func TestCreateAccount(t *testing.T) {
	store := &fakeAccountStore{
		createFunc: func(
			ctx context.Context,
		) (*domain.Account, error) {
			return &domain.Account{
				ID:             42,
				BalanceKopecks: 0,
			}, nil
		},
	}

	app := &application{
		accounts: store,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		nil,
	)

	response := httptest.NewRecorder()

	app.newRouter().ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusCreated,
			response.Code,
		)
	}

	expected := `{"id":42,"balance_kopecks":0}`

	if strings.TrimSpace(
		response.Body.String(),
	) != expected {
		t.Errorf(
			"ожидался ответ %s, получен %s",
			expected,
			response.Body.String(),
		)
	}
}

func TestGetAccountNotFound(t *testing.T) {
	store := &fakeAccountStore{
		getByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*domain.Account, error) {
			return nil, domain.ErrAccountNotFound
		},
	}

	app := &application{
		accounts: store,
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/accounts/999",
		nil,
	)

	response := httptest.NewRecorder()

	app.newRouter().ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusNotFound,
			response.Code,
		)
	}

	expected := `"code":"ACCOUNT_NOT_FOUND"`

	if !strings.Contains(
		response.Body.String(),
		expected,
	) {
		t.Errorf(
			"ответ не содержит %s: %s",
			expected,
			response.Body.String(),
		)
	}
}

func TestDepositInvalidAmount(t *testing.T) {
	store := &fakeAccountStore{
		depositFunc: func(
			ctx context.Context,
			accountID int64,
			amount int64,
		) (*domain.Account, error) {
			return nil, domain.ErrInvalidAmount
		},
	}

	app := &application{
		accounts: store,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/accounts/1/deposit",
		strings.NewReader(
			`{"amount_kopecks":0}`,
		),
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

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if !strings.Contains(
		response.Body.String(),
		`"code":"INVALID_AMOUNT"`,
	) {
		t.Errorf(
			"получен неожиданный ответ: %s",
			response.Body.String(),
		)
	}
}
