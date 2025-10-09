package main

import (
	"tcp-server.com/m/internal/config"
	"tcp-server.com/m/internal/server"
)

func main() {
	server := server.NewServer(config.Port)
	server.Start()
}
