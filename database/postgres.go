package database

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}