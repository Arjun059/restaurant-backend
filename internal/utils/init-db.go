package utils

import (
	_ "modernc.org/sqlite"
	"gorm.io/gorm"
	"os"
	"errors"
	"gorm.io/driver/postgres"
	"fmt"
)

func InitDB() (*gorm.DB, error) {
 	APP_ENV := os.Getenv("APP_ENV");
	var db *gorm.DB
	var err error
		
	// for production development use sqlLite db
	dsn := os.Getenv("POSTGRESS_SQL_URL")
	if dsn == "" {
			return nil, errors.New("POSTGRESS_SQL_URL environment variable is not set")
	}
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
			return nil, fmt.Errorf("failed to connect to Supabase: %w", err)
	}

	if err != nil {
		return nil, err
	}

	if APP_ENV != "prod" {
		fmt.Println("Development DB connected successfully")
	} else {
		fmt.Println("Production DB connected successfully")
	}

	return db, nil
}
