package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"rent/internal/config"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.GetDBConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database")
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, err
}
