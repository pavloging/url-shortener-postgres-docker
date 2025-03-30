package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Подключение к новой БД shortener
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	// Создание таблиц в новой БД
	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS url (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		log.Fatalf("Cannot create table: %v\n", err)
	}

	// Умная вставка с обработкой дубликатов
	var insertedID int
	err = db.QueryRow(
		context.Background(),
		`INSERT INTO url(alias, url) VALUES($1, $2)
		 ON CONFLICT (alias) DO UPDATE SET url = EXCLUDED.url
		 RETURNING id`,
		"example123r",
		"https://example12r345.com",
	).Scan(&insertedID)

	if err != nil {
		log.Fatalf("Insert failed: %v\n", err)
	}

	fmt.Printf("Operation completed successfully! ID: %d\n", insertedID)
}
