package database

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq" // postgres driver
)

var DB *sql.DB

func ConnectPostgres() *sql.DB {
    if DB != nil {
        return DB
    }

    connStr := "host=localhost port=5432 user=postgres password=123456 dbname=prestasi_db sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Gagal buka koneksi Postgres:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Gagal koneksi Postgres:", err)
    }

    log.Println("PostgreSQL berhasil terkoneksi")
    DB = db
    return DB
}
