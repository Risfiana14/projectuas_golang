// database/postgres.go
package database

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {
    // HARD CODE — INI YANG 100% JALAN DI SEMUA LAPTOP KAMPUS
    connStr := "host=localhost port=5432 user=postgres password=123456 dbname=prestasi_db sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Gagal buka PostgreSQL:", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal("PostgreSQL tidak respons:", err)
    }

    log.Println("PostgreSQL berhasil terkoneksi (password: 123456)")
    return db
}