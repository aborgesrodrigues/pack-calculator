package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"pack-calculator/internal/common"
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

func (h *Handler) SavePackSize(w http.ResponseWriter, r *http.Request) {
	packSizeBatch := &common.PackSizeBatch{}
	if err := json.NewDecoder(r.Body).Decode(packSizeBatch); err != nil {
		h.logger.Error("Unable to decode request body in handler.SavePackSize", "error", err)
		writeResponse(w, http.StatusBadRequest, "Payload in the wrong format")
		return
	}

	if len(packSizeBatch.Sizes) == 0 {
		h.logger.Error("No sizes passed", "sizes", packSizeBatch.Sizes)
		if err := writeResponse(w, http.StatusBadRequest, "no sizes passed"); err != nil {
			h.logger.Error(err.Error())
		}
		return
	}

	if err := h.service.SavePackSize(r.Context(), packSizeBatch); err != nil {
		h.logger.Error("Error saving pack sizes", "error", err)

		if err := writeResponse(w, http.StatusInternalServerError, "error saving pack sizes"); err != nil {
			h.logger.Error(err.Error())
		}
		return
	}

	if err := writeResponse(w, http.StatusOK, "Pack sizes saved"); err != nil {
		h.logger.Error(err.Error())
	}
}
