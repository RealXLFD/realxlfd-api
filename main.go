package main

import (
	_ "github.com/mattn/go-sqlite3"
	"rpics-docker/server/db"
)

func main() {
	sql := db.Connect()

}
