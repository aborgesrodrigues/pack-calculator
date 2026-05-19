package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"pack-calculator/internal/service"
)

type Handler struct {
	logger  *slog.Logger
	service service.SVCInterface
}

func NewHandler(logger *slog.Logger) (*Handler, error) {
	svc, err := service.NewService(logger)

	if err != nil {
		return nil, fmt.Errorf("error instantiating service: %w", err)
	}

	return &Handler{
		logger:  logger,
		service: svc,
	}, nil
}

func writeResponse(w http.ResponseWriter, status int, message any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(message); err != nil {
		return fmt.Errorf("error writing response %w", err)
	}

	return nil
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "ok")
}
