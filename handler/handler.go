package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"strongdm/counter"
)

// Handler holds the rate limiting counter and provides HTTP request handling
type Handler struct {
	counter *counter.Counter
}

// New creates a new HTTP handler with rate limiting
func New() *Handler {
	return &Handler{
		counter: counter.New(),
	}
}

// HandleRequest processes HTTP requests with rate limiting
func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteHost = r.RemoteAddr
	}

	const limitPerMinute = 120

	info := h.counter.Add(remoteHost, limitPerMinute, 1)

	w.Header().Set("Content-Type", "application/json")

	if info.Allowed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}

	jsonData, _ := json.MarshalIndent(info, "", "  ")
	_, _ = w.Write(jsonData)
}
