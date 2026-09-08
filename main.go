package main

import (
	"log"
	"net/http"
	"strings"
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

func respond(transcript string) string {
	t := strings.ToLower(transcript)
	switch {
	case strings.Contains(t, "hello") || strings.Contains(t, "hi"):
		return "Good day to you. How may I be of service?"
	case strings.Contains(t, "time"):
		return "I'm afraid I don't yet have a clock at hand, but I shall acquire one shortly."
	case strings.Contains(t, "thank"):
		return "It is entirely my pleasure."
	default:
		return "I heard you say: \"" + transcript + "\". I'm still learning how best to assist with that."
	}
}

