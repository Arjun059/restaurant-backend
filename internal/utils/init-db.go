package utils

import (
	"gorm.io/driver/sqlite"
	_ "modernc.org/sqlite"
	"gorm.io/gorm"
	"os"
	"errors"
	"gorm.io/driver/postgres"
	"fmt"
)

func InitDB() (*gorm.DB, error) {
	// for local development use sqlLite db
 	APP_ENV := os.Getenv("APP_ENV");
	var db *gorm.DB
	var err error
	if APP_ENV == "production" {
			// for production development use sqlLite db
			dsn := os.Getenv("POSTGRESS_SQL_URL")
			if dsn == "" {
					return nil, errors.New("POSTGRESS_SQL_URL environment variable is not set")
			}
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
					return nil, fmt.Errorf("failed to connect to Supabase: %w", err)
			}
	} else {
     // for development
			db, err = gorm.Open(sqlite.Dialector{
				DriverName: "sqlite", // Force GORM to use modernc SQLite
				DSN:        "app.db", // Database file
			}, &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}
	return db, nil
}
