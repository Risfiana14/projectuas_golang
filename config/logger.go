// config/logger.go (simple, adaptasi dari modul)
package config

import (
    "log"
)

func InitLogger() {
    log.Println("Logger initialized")  // Bisa extend untuk rotating files
}