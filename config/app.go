package config

import (
    "log"
    "projectuas/database"
    "projectuas/route"
    "projectuas/app/repository"

    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
)

// LoadEnv — satu-satunya fungsi LoadEnv di seluruh project
func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Println("File .env tidak ditemukan, menggunakan environment variable langsung")
    }
}

// NewApp — dipanggil di main.go
func NewApp() *fiber.App {
    app := fiber.New()

    // Middleware logger sederhana (sesuai modul)
    app.Use(func(c *fiber.Ctx) error {
        log.Printf("Request: %s %s", c.Method(), c.Path())
        return c.Next()
    })

    // Koneksi Mongodb
    mongoClient := database.ConnectMongo()
    repository.InitMongo(mongoClient)

    // Koneksi Postgres
    pg := database.ConnectPostgres()
    repository.Init(pg) // Init repository dengan Postgres

    // Register semua route
    route.Setup(app)

    return app
}