package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON sends v as a JSON body with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError sends the project's uniform error shape: {"error": message}.
// Messages are human-readable and never include secrets or internal details
// (AC-15, AC-16).
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
