package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	return sql.Open("mysql", cfg.FormatDSN())
}

func InitStorage(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		log.Fatal("error in connecting to db")
	}

	log.Println("connected to db")
	return nil
}
