package main

import (
	"errors"
	"net"
	"net/url"
	"os"
)

func databaseURL() (string, error) {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn, nil
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	if host == "" ||
		port == "" ||
		database == "" ||
		user == "" ||
		password == "" {

		return "", errors.New(
			"параметры подключения к PostgreSQL не заданы",
		)
	}

	dsn := &url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			user,
			password,
		),
		Host: net.JoinHostPort(
			host,
			port,
		),
		Path: database,
	}

	query := dsn.Query()
	query.Set("sslmode", "disable")

	dsn.RawQuery = query.Encode()

	return dsn.String(), nil
}
