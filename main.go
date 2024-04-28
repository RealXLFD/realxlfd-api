package main

import (
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	for {
		time.Sleep(time.Second)

		println(strconv.Atoi(result))
	}
}
