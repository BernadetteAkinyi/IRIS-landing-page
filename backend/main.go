package main

import (
	"log"
	"net/http"
	"os"

	"backend/assistant"
	"backend/handlers"
	"backend/store"
)

func main() {
	// Initialize store and assistant engine
	dataStore := store.New()
	assistantEngine := assistant.NewEngine(dataStore)
	api := handlers.NewAPI(dataStore, assistantEngine)

	mux := http.NewServeMux()

	// Assistant route
	mux.HandleFunc("/api/assistant", api.AssistantHandler)

	// Accessibility preferences routes
	mux.HandleFunc("/api/preferences", api.PreferencesHandler)

	// User management routes
	mux.HandleFunc("/api/users/login", api.LoginHandler)
	mux.HandleFunc("/api/users/register", api.RegisterHandler)

	// Services & Events routes
	mux.HandleFunc("/api/services", api.ServicesHandler)
	mux.HandleFunc("/api/events", api.EventsHandler)

	// Health check route
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		handlers.EnableCORS(w)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"iris-butler-backend"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("IRIS Butler backend server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
