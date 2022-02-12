package service

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthService is a HTTP Consumer for health checking
type HealthService struct {
	logger *log.Logger
}

// NewHealth creates a new Health handler
func NewHealth(l *log.Logger) *HealthService {
	return &HealthService{l}
}

// ServeHTTP implements the handler interface
func (h *HealthService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("ok")
	h.logger.Println(err)
}
