package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl != "" {
		fmt.Printf("database url: %v", databaseUrl)
		db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})

		if err != nil {

			return nil, fmt.Errorf("failed to connect to the database: %w", err)
		}

		fmt.Println("Connected to Railway database")
		return db, nil

	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
		)

		localdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {

			return nil, fmt.Errorf("failed to connect to the database: %w", err)
		}

		fmt.Println("Connected to local database")
		return localdb, nil
	}

}
