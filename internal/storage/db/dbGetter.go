package db

import (
	"database/sql"
	"log"
	"rent/internal/config"
)

var db *sql.DB = nil

func UpdateDB() error {
	newDb, err := Connect(config.GetSingletonConfig())

	if err != nil {
		log.Fatal("Failed to connect to database")
		return err
	}

	db = newDb

	return nil
}

func GetSingletonDB() (*sql.DB, error) {
	if db == nil {
		err := UpdateDB()

		if err != nil {
			log.Fatal("Failed update db")
			return nil, err
		}
	}

	return db, nil
}
