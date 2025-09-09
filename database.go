package main

import (
	"database/sql"

	//_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite", "./cats.db")
	if err != nil {
		return err
	}
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS cat (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        age INTEGER
    );`

	_, err = db.Exec(createTableSQL)
	return err
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
