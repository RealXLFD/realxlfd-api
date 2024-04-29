package main

import (
	"rpics-docker/server"
	"rpics-docker/server/func/rpic"
	"rpics-docker/server/func/welcome"
)

func main() {
	port := readConfig()
	server.Run(port, rpic.Serve, welcome.Serve)
}
