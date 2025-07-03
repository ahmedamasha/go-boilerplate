package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

// GetDB returns a singleton instance of the database connection
func GetDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "microservice_db"),
			getEnv("DB_SSL_MODE", "disable"),
		)

		log.Printf("Connecting to database with connection string: %s", connStr)

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error opening database connection: %v", err)
			return
		}

		// Set connection pool settings
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		// Test the connection
		err = db.Ping()
		if err != nil {
			log.Printf("Error connecting to database: %v", err)
			return
		}

		log.Printf("Successfully connected to database")
	})

	return db, err
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func NewDatabaseConnection(config *Config) (*sql.DB, error) {
	dsn := config.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database: %s:%s/%s",
		config.Database.Host, config.Database.Port, config.Database.DBName)

	return db, nil
}
