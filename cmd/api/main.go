package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"svipp-server/internal/config"
	"svipp-server/internal/version"
)

type server struct {
	config   *config.Config
	services *config.Services
}

func main() {
	err := run()
	if err != nil {
		trace := string(debug.Stack())
		log.Printf("Error: %s\nTrace: %s", err.Error(), trace)
		os.Exit(1)
	}
}

func run() error {

	fmt.Printf("\nRunning svipp-server version: %s\n", version.Get())

	cfg, services, err := config.New()
	if err != nil {
		return err
	}

	srv := &server{
		config:   cfg,
		services: services,
	}

	return srv.serveHTTP()
}
