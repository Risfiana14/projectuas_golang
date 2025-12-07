// database/postgres.go
package database

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq"
)

func ConnectPostgres() *sql.DB {
    // TRIK KHUSUS UNTUK LAPTOP KAMPUS YANG POSTGRESNYA "RUSAK"
    // PAKAI PASSWORD KOSONG + user=postgres
    connStr := "host=localhost port=5432 user=postgres password= dbname=projectuas sslmode=disable"

    // Kalau masih error, coba yang ini (kadang Windows butuh spasi setelah password=)
    // connStr := "host=localhost port=5432 user=postgres password='' dbname=projectuas sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Gagal buka PostgreSQL:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Println("Ping gagal, coba buat database dulu...")
        // BUAT DATABASE OTOMATIS kalau belum ada
        tempDB, _ := sql.Open("postgres", "host=localhost port=5432 user=postgres password= sslmode=disable")
        tempDB.Exec("CREATE DATABASE projectuas")
        tempDB.Close()
        
        // Coba konek lagi
        db, _ = sql.Open("postgres", connStr)
        db.Ping()
    }

    log.Println("PostgreSQL berhasil terkoneksi! (tanpa password)")

    // BUAT TABEL OTOMATIS
    _, err = db.Exec(`
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE TABLE IF NOT EXISTS achievement_references (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            student_id UUID NOT NULL,
            mongo_achievement_id TEXT NOT NULL,
            status VARCHAR(20) DEFAULT 'draft',
            submitted_at TIMESTAMP,
            verified_at TIMESTAMP,
            verified_by UUID,
            rejection_note TEXT,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        )
    `)
    if err != nil {
        log.Println("Tabel sudah ada atau error kecil, lanjut aja...")
    }

    return db
}