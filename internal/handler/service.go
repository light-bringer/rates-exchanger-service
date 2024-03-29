package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// Handler represents an HTTP handler for exchange rates.

// GetLatestRates handles requests for the latest exchange rates.
func (h *Handler) GetLatestRates(w http.ResponseWriter, r *http.Request) {
	limit := uint64(0)
	var err error
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			slog.Error("Invalid limit", "error", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}

	rates, err := h.service.FetchLatestExchangeRates(limit)
	if err != nil {
		slog.Error("Failed to fetch latest exchange rates", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"base":  baseCurrency,
		"rates": rates,
	}

	w.Header().Set(contentTypeHeader, contentType)
	json.NewEncoder(w).Encode(response)
}

// HealthCheck handles requests for the health check.
func (h *Handler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// GetExchangeRate handles requests for the exchange rate for a specific date.
func (h *Handler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("calculationDay")

	slog.Debug("fetching exchange rate for date", "date", date)

	// Validate date format (YYYY-MM-DD)
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		slog.Error("Invalid date format", "error", err)
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	limit := uint64(0)
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			slog.Error("Invalid limit", "error", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}

	rates, err := h.service.FetchRatesForDate(date, limit)
	if err != nil {
		slog.Error("Failed to fetch latest exchange rates", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"base":  baseCurrency,
		"rates": rates,
	}

	w.Header().Set(contentTypeHeader, contentType)
	json.NewEncoder(w).Encode(response)
}

// GetStatistics handles requests for the exchange rate statistics.
func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	slog.Debug("logging statistics for exchange rates", "request", r.Header)

	days := uint64(defaultRange)
	var err error
	dayStr := r.URL.Query().Get("range")
	if dayStr != "" {
		days, err = strconv.ParseUint(dayStr, 10, 64)
		if err != nil {
			slog.Error("Invalid limit", "error", err)
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}

	stats, err := h.service.GetRateStatistics(days)
	if err != nil {
		slog.Error("Failed to fetch rate statistics", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	detailedStats := make(map[string]interface{})
	for currency, stat := range stats {
		detailedStats[currency] = map[string]interface{}{
			"average": stat.AvgRate,
			"min":     stat.MinRate,
			"max":     stat.MaxRate,
		}
	}

	response := map[string]interface{}{
		"base":          baseCurrency,
		"rates_analyze": detailedStats,
	}

	w.Header().Set(contentTypeHeader, contentType)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/rates/latest", h.GetLatestRates)
	mux.HandleFunc("/rates/{calculationDay}", h.GetExchangeRate)
	mux.HandleFunc("/rates/analyze", h.GetStatistics)
	mux.HandleFunc("/health", h.HealthCheck)
	return mux
}
