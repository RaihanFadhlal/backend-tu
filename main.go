package main

import (
	"log"
	"backendtku/app/server"
	"backendtku/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	s := server.NewServer(cfg)
	s.Run()
}
