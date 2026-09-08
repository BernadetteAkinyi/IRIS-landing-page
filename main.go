package main

import (
	"log"
	"net/http"
)

// Request body sent from the browser
type AssistantRequest struct {
	Transcript string `json:"transcript"`
}

// Response sent back to the browser
type AssistantResponse struct {
	Reply string `json:"reply"`
}

func main() {
	log.Println("Butler backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}