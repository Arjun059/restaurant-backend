package utils

import (
	"gorm.io/driver/sqlite"
	_ "modernc.org/sqlite"
	"gorm.io/gorm"
	// "os"
	// "gorm.io/driver/postgres"
	// "log"
)

func InitDB() (*gorm.DB, error) {
	// for local development use sqlLite db
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite", // Force GORM to use modernc SQLite
		DSN:        "app.db", // Database file
	}, &gorm.Config{})

	// for production development use sqlLite db
	// dsn := os.Getenv("POSTGRESS_SQL_URL")
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 		log.Fatalf("failed to connect to Supabase: %v", err)
	// }

	if err != nil {
		return nil, err
	}
	return db, nil
}
