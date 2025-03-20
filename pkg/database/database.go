package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ispindel.piwo.org/internal/models"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_USER", "ispindel"),
		getEnvOrDefault("DB_PASSWORD", "Kochanapysia1"),
		getEnvOrDefault("DB_NAME", "ispindel"),
		getEnvOrDefault("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Nie udało się połączyć z bazą danych:", err)
	}

	// Automatyczna migracja schematu
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Nie udało się zmigrować schematu bazy danych:", err)
	}

	DB = db
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 