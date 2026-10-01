package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ars-Mik/bank-transactions-api/internal/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	db       *sql.DB
	accounts *repository.AccountRepository
}

func (app *application) newRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /ready", app.readinessHandler)

	mux.HandleFunc("POST /accounts", app.createAccountHandler)
	mux.HandleFunc("GET /accounts", app.listAccountsHandler)
	mux.HandleFunc("GET /accounts/{id}", app.getAccountHandler)
	mux.HandleFunc("POST /accounts/{id}/deposit", app.depositHandler)
	mux.HandleFunc("POST /transfers", app.transferHandler)
	mux.HandleFunc("GET /accounts/{id}/transactions", app.accountTransactionsHandler)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ok",
		},
	)
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
		db:       db,
		accounts: accountRepository,
	}

	server := &http.Server{
		Addr:              ":8080",
		Handler:           app.newRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	shutdownSignal, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		log.Println(
			"Bank Transactions API запущен: http://localhost:8080",
		)

		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			log.Fatal(
				"HTTP-сервер завершился с ошибкой: ",
				err,
			)
		}

	case <-shutdownSignal.Done():
		log.Println(
			"Получен сигнал завершения. Останавливаем сервер...",
		)
	}

	shutdownContext, cancelShutdown := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf(
			"Не удалось корректно завершить HTTP-сервер: %v",
			err,
		)
	} else {
		log.Println(
			"HTTP-сервер корректно остановлен",
		)
	}
}
