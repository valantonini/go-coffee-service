package health

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthService is a HTTP Consumer for health checking
type HealthService struct {
	logger *log.Logger
}

// NewHealthService creates a new Health handler
func NewHealthService(l *log.Logger) *HealthService {
	return &HealthService{l}
}

// ServeHTTP implements the handler interface
func (h *HealthService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode("ok")
	if err != nil {
		h.logger.Println(err)
	}
}
