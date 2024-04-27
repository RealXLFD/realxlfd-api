package main

import (
	"github.com/davidbyttow/govips/v2/vips"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()
}
