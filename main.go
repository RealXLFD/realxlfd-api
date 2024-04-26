package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	row := db.QueryRow("SELECT * FROM Images;")
	hash := ""
	err = row.Scan(&hash)
	if err != nil {
		panic(err)
	}
	print(hash)
}
