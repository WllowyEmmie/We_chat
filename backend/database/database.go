package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	dsn := os.Getenv("DATABASE_URL")
	
	if dsn == "" {
		fmt.Println("DATABASE_URL not found, falling back to individual env vars...")

		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASS")
		name := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		// Railway usually requires sslmode=require, local can use disable
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable" // change to "require" if you want
		}

		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, pass, host, port, name, sslmode,
		)
	}

	fmt.Println("🔌 Connecting to DB:", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil

}
