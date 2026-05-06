package metrics

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.Handler that serves a JSON snapshot of the registry.
// GET /metrics  →  { "counter_name": 42, ... }
func Handler(r *Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := r.Snapshot()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(snap); err != nil {
			// Headers already sent; nothing useful we can do.
			_ = err
		}
	})
}
