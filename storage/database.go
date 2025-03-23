package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type AbsDatabase interface {
	ConnectToDB()
	Insert(tableName string, text string)
	Close()
}

type Database struct {
	db *sql.DB
}

func NewDatabase() AbsDatabase {
	return &Database{}
}

func (d *Database) ConnectToDB() {
	connStr := fmt.Sprintf(
        "host=database port=5432 user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_NAME"),
    )
	db, err := sql.Open("postgres", connStr)
	failOnError(err, "Failed to connect to database")

	err = db.Ping()
	failOnError(err, "Failed to connect to database")

	d.db = db
}

func (db *Database) Insert(tableName string, text string) {
	query := fmt.Sprintf("INSERT INTO %s (message) VALUES ($1)", tableName)
	_, err := db.db.Exec(query, text)
	failOnError(err, "Failed to add record to database")
}

func (db *Database) Close() {
	db.db.Close()
}