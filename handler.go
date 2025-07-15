package main

import (
	"encoding/json"
	"net"
	"net/http"
)

var counter = NewCounter()

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteHost = r.RemoteAddr
	}

	const limitPerMinute = 120

	info := counter.Add(remoteHost, limitPerMinute, 1)

	w.Header().Set("Content-Type", "application/json")

	if info.Allowed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}

	jsonData, _ := json.MarshalIndent(info, "", "  ")
	w.Write(jsonData)
}