package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// Response struct to hold the UUID and server info
type Response struct {
	UUID      string `json:"uuid"`
	Hostname  string `json:"hostname"`
	Timestamp string `json:"timestamp"`
 	ClientIP  string `json:"clientip"`
}

// HealthCheckResponse struct to hold the health check status
type HealthCheckResponse struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
}

// Handler to return a random UUID and set a cookie
func uuidHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a new UUID
	id := uuid.New()

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "could not get hostname", http.StatusInternalServerError)
		return
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response with the UUID and hostname
	response := Response{
		UUID:      id.String(),
		Hostname:  hostname,
		Timestamp: time.Now().In(loc).Format("2006-01-02 15:04:05"),
  		ClientIP:  r.RemoteAddr,
	}

	// Marshal the response into JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "uuid",
		Value: id.String(),
		Path:  "/",
	})

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Health check handler
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		http.Error(w, "could not get hostname", http.StatusInternalServerError)
		return
	}

	// Create a response with the status and hostname
	response := HealthCheckResponse{
		Status:   "ok",
		Hostname: hostname,
	}

	// Marshal the response into JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	// Set up the HTTP handlers
	http.HandleFunc("/uuid", uuidHandler)
	http.HandleFunc("/healthcheck", healthCheckHandler)

	// Start the server
	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
