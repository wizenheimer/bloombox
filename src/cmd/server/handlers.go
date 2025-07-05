package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/pkg/emailchecker"
)

// handleIndex handles the root endpoint and returns API documentation
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": "Email Validation API",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"POST /validate":        "Validate single email address",
			"POST /batch":           "Validate multiple email addresses",
			"GET /validators":       "List available validators",
			"PUT /validators/:name": "Enable/disable specific validator",
			"GET /health":           "Health check endpoint",
		},
		"validators": s.checker.GetValidators(),
	})
}

// handleValidate handles single email validation requests
func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req emailchecker.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	var result *emailchecker.CheckResult
	if len(req.Validators) > 0 {
		result = s.checker.CheckWithValidators(req.Email, req.Validators)
	} else {
		result = s.checker.Check(req.Email)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleBatch handles batch email validation requests
func (s *Server) handleBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Emails     []string `json:"emails"`
		Validators []string `json:"validators,omitempty"`
		Timeout    int      `json:"timeout,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Emails) == 0 {
		http.Error(w, "At least one email is required", http.StatusBadRequest)
		return
	}

	if len(req.Emails) > 100 {
		http.Error(w, "Maximum 100 emails per batch", http.StatusBadRequest)
		return
	}

	results := make([]*emailchecker.CheckResult, len(req.Emails))

	// Process emails concurrently
	type result struct {
		index int
		data  *emailchecker.CheckResult
	}

	resultChan := make(chan result, len(req.Emails))

	for i, email := range req.Emails {
		go func(idx int, em string) {
			var res *emailchecker.CheckResult
			if len(req.Validators) > 0 {
				res = s.checker.CheckWithValidators(em, req.Validators)
			} else {
				res = s.checker.Check(em)
			}
			resultChan <- result{index: idx, data: res}
		}(i, email)
	}

	// Collect results
	for i := 0; i < len(req.Emails); i++ {
		res := <-resultChan
		results[res.index] = res.data
	}

	response := map[string]interface{}{
		"results": results,
		"summary": map[string]interface{}{
			"total":      len(results),
			"valid":      countValid(results),
			"invalid":    len(results) - countValid(results),
			"disposable": countDisposable(results),
			"free":       countFree(results),
			"role":       countRole(results),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleValidators handles validator management requests
func (s *Server) handleValidators(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		validators := s.checker.GetValidators()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"validators": validators,
		})

	case http.MethodPut:
		// Parse validator name from URL path
		path := strings.TrimPrefix(r.URL.Path, "/validators/")
		if path == "" {
			http.Error(w, "Validator name is required", http.StatusBadRequest)
			return
		}

		var req struct {
			Enabled bool `json:"enabled"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := s.checker.SetValidatorEnabled(path, req.Enabled); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"validator": path,
			"enabled":   req.Enabled,
			"message": fmt.Sprintf("Validator %s %s", path,
				func() string {
					if req.Enabled {
						return "enabled"
					} else {
						return "disabled"
					}
				}()),
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	validators := s.checker.GetValidators()
	enabledCount := 0
	for _, enabled := range validators {
		if enabled {
			enabledCount++
		}
	}

	health := map[string]interface{}{
		"status":             "healthy",
		"timestamp":          time.Now(),
		"enabled_validators": enabledCount,
		"total_validators":   len(validators),
		"cache_size":         s.config.CacheSize,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
