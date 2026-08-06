package handlers

import (
	"net/http"
)

func TelemetryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET method is supported")
		return
	}

	metrics := map[string]string{
		"experienceYears": "5+",
		"dailyEvents":     "5M+",
		"avgLatencyMs":    "<50ms",
		"activeSessions":  "3,000+",
		"systemStatus":    "OPERATIONAL",
	}

	respondJSON(w, http.StatusOK, metrics)
}
