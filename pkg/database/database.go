package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := getEnvOrDefault("DB_HOST", "mysql18.mydevil.net")
	user := getEnvOrDefault("DB_USER", "m1270_ispindel")
	password := mustGetEnv("DB_PASSWORD")
	dbname := getEnvOrDefault("DB_NAME", "m1270_ispindel")
	port := getEnvOrDefault("DB_PORT", "3306")

	log.Printf("Próba połączenia z bazą danych MySQL: %s@%s:%s/%s", user, host, port, dbname)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Nie udało się połączyć z bazą danych MySQL:", err)
	}

	log.Println("Połączenie z bazą danych MySQL udane!")

	DB = db

	// Uruchom migracje
	RunMigrations()
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Wymagana zmienna środowiskowa %s nie jest ustawiona", key)
	}
	return v
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
