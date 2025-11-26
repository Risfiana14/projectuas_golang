package config

import (
    "projectuas/database"
    "projectuas/middleware"
    "projectuas/route"

    "github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
    app := fiber.New()
    app.Use(middleware.Logger()) // ← SEKARANG SUDAH BENAR!

    pgDB := database.ConnectPostgres()
    mongoClient := database.ConnectMongo()

    route.Setup(app, pgDB, mongoClient)
    return app
}