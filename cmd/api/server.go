package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	defaultIdleTimeout    = 2 * time.Minute
	defaultReadTimeout    = 10 * time.Second
	defaultWriteTimeout   = 30 * time.Second
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
		s.config.DB.DBPool.Close()
		log.Printf("Server exited: %v", err)
		return error(errors.New("server exited"))
	}

	// This line will only be reached if ListenAndServe returns without error (which is unlikely)
	s.config.DB.DBPool.Close()
	log.Printf("Server exited")
	return error(errors.New("server exited"))
}
