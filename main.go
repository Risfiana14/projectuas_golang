// main.go
package main

import (
    "log"
    "os"
    "projectuas/config"
)

func main() {
    config.LoadEnv()
    app := config.NewApp()  // SEKARANG SUDAH ADA!
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "3000"
    }
    log.Printf("Server jalan di http://localhost:%s", port)
    log.Fatal(app.Listen(":" + port))
}