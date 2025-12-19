// main.go
package main

import (
	"log"
	"os"
	"projectuas/config"

	_ "projectuas/docs" // ⬅️ WAJIB
)

// @title Achievement Management API
// @version 1.0
// @description API Sistem Manajemen Prestasi Mahasiswa
// @termsOfService http://swagger.io/terms/

// @contact.name Kelompok UAS PBL
// @contact.email pbl@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	config.LoadEnv()
	app := config.NewApp()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server jalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
