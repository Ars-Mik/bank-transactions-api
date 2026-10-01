package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ars-Mik/bank-transactions-api/internal/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	accounts *repository.AccountRepository
}

func (app *application) newRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("POST /accounts", app.createAccountHandler)
	mux.HandleFunc("GET /accounts", app.listAccountsHandler)
	mux.HandleFunc("GET /accounts/{id}", app.getAccountHandler)
	mux.HandleFunc("POST /accounts/{id}/deposit", app.depositHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func main() {

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal("failed to initialize database: ", err)
	}

	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(
		30 * time.Minute,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	err = db.PingContext(ctx)
	cancel()

	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}

	log.Println("Подключение к PostgreSQL успешно установлено")

	accountRepository := repository.NewAccountRepository(db)

	app := &application{
		accounts: accountRepository,
	}

	server := &http.Server{
		Addr:              ":8080",
		Handler:           app.newRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Bank Transactions API запущен: http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
