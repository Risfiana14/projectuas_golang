// config/app.go
package config

import (
    "log"
    "projectuas/database"
    "projectuas/route"

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

    // Koneksi database
    pgDB := database.ConnectPostgres()
    mongoClient := database.ConnectMongo()

    // Register semua route
    route.Setup(app, pgDB, mongoClient)

    return app
}