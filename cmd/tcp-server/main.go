package main

import "tcp-server.com/m/internal/server"

func main() {
	server := server.NewServer(":3000")
	server.Start()
}
