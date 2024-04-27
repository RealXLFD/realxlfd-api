package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	println(uint(-100 % 3))
}
