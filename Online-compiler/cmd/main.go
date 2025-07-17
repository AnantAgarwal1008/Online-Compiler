package main

import (
	"log"
	"net/http"

	"github.com/AnantAgarwal1008/online-compiler/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/run", handler.RunCodeHandler)

	// Wrap the mux with CORS support
	withCORS := func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (for local testing)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// For other requests, forward to mux
		mux.ServeHTTP(w, r)
	}

	log.Println("Server running at http://localhost:8089")
	if err := http.ListenAndServe(":8089", http.HandlerFunc(withCORS)); err != nil {
		log.Fatal("Server failed:", err)
	}
}
