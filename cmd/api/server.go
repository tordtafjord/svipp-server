package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (s *server) serveHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.HTTPPort),
		Handler:      s.routes(),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	log.Printf("Server is starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}

	// This line will only be reached if ListenAndServe returns without error (which is unlikely)
	s.config.DB.DBPool.Close()
	log.Println("Server exited")
	return error(errors.New("Server exited"))
}
