package main

import (
	"rpics-docker/server"
)

func main() {
	server.ServeWelcome()
	server.ServeRpic()
	server.Run()
}
