package service

import (
	"fmt"
	"log"
	"net/http"
)

// HealthService is a HTTP Handler for health checking
type HealthService struct {
	logger *log.Logger
}

// NewHealth creates a new Health handler
func NewHealth(l *log.Logger) *HealthService {
	return &HealthService{l}
}

// ServeHTTP implements the handler interface
func (h *HealthService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", "\"ok\"")
}
