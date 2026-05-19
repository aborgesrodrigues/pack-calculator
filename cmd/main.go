package main

import (
	"net/http"
	"pack-calculator/cmd/handler"
	"pack-calculator/internal/utils"
	"time"
)

func main() {
	logger := utils.NewLogger()

	handler, err := handler.NewHandler(logger)
	if err != nil {
		logger.Error("error instantiating handler", "error", err)
		panic("error instantiating handler")
	}

	// add routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  // TODO: needs to be changed?
		WriteTimeout: 5 * time.Second,  // TODO: needs to be changed?
		IdleTimeout:  15 * time.Second, // TODO: needs to be changed?
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("could not start the server, error", "error", err)
		panic("could not start the server")
	}
}
