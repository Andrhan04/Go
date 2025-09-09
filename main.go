package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./Lesson.db")
	if err != nil {
		print("AAAAAAAAAAAAAAA\n\n\n")
		panic(err)
	}
	print("AAAAAAAAAAAAAAA\n\n\n")
	state, _ := db.Prepare("CREATE TABLE cat (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, age INTEGER);")
	state.Exec()

}
