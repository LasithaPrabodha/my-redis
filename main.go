package main

import (
	"github.com/LasithaPrabodha/redis-like-server/cmd/server"
)

func main() {
	server.Start(":6379")
}
