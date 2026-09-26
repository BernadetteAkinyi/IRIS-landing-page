package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type AssistantRequest struct {
	Transcript string `json:"transcript"`
}

type ButlerAction struct {
	Type    string `json:"type"`
	Payload string `json:"payload,omitempty"`
}

type AssistantResponse struct {
	Reply  string        `json:"reply"`
	Action *ButlerAction `json:"action,omitempty"`
}

func parseCommand(transcript string) (string, *ButlerAction) {
	t := strings.ToLower(transcript)

	switch {
	case strings.Contains(t, "baby pink") || strings.Contains(t, "pink"):
		return "Right away. I have changed the page to baby pink.", &ButlerAction{
			Type:    "SET_THEME",
			Payload: "baby-pink",
		}

	case strings.Contains(t, "easier to read") || strings.Contains(t, "easier for me to read") || strings.Contains(t, "reading mode") || strings.Contains(t, "dyslexia"):
		return "Certainly. I have applied reading mode with clear typography and relaxed spacing.", &ButlerAction{
			Type:    "SET_READING_MODE",
			Payload: "true",
		}

	case (strings.Contains(t, "increase") && strings.Contains(t, "text")) || strings.Contains(t, "bigger") || strings.Contains(t, "larger") || strings.Contains(t, "increase font"):
		return "At once. I have increased the text size.", &ButlerAction{
			Type:    "SET_TEXT_SIZE",
			Payload: "large",
		}

	case (strings.Contains(t, "decrease") && strings.Contains(t, "text")) || strings.Contains(t, "smaller") || strings.Contains(t, "decrease font"):
		return "Certainly. I have decreased the text size.", &ButlerAction{
			Type:    "SET_TEXT_SIZE",
			Payload: "small",
		}

	case strings.Contains(t, "text to speech") || strings.Contains(t, "text-to-speech") || strings.Contains(t, "read"):
		return "Text-to-speech is enabled. Reading page content for you.", &ButlerAction{
			Type:    "READ_PAGE",
			Payload: "true",
		}

	case strings.Contains(t, "simpler") || strings.Contains(t, "simplify"):
		return "I have simplified the interface to highlight essential content.", &ButlerAction{
			Type:    "SIMPLIFY",
			Payload: "true",
		}

	case strings.Contains(t, "high contrast") || strings.Contains(t, "dark mode"):
		return "High contrast has been enabled for clearer visibility.", &ButlerAction{
			Type:    "HIGH_CONTRAST",
			Payload: "true",
		}

	case strings.Contains(t, "reduce motion") || strings.Contains(t, "stop animation"):
		return "Reduced motion is enabled. Animations have been minimized.", &ButlerAction{
			Type:    "REDUCE_MOTION",
			Payload: "true",
		}

	case strings.Contains(t, "reset") || strings.Contains(t, "default"):
		return "I have restored the default interface.", &ButlerAction{
			Type:    "RESET",
			Payload: "true",
		}

	case strings.Contains(t, "hello") || strings.Contains(t, "hi") || strings.Contains(t, "hey"):
		return "Good day. I am IRIS Butler. How may I assist you with your accessibility preferences?", nil

	case strings.Contains(t, "thank"):
		return "It is entirely my pleasure.", nil

	default:
		return "I heard: \"" + transcript + "\". You can ask me to change the page to baby pink, make text easier to read, increase text size, or simplify the page.", nil
	}
}

func assistantHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AssistantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reply, action := parseCommand(req.Transcript)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AssistantResponse{
		Reply:  reply,
		Action: action,
	})
}

func main() {
	http.HandleFunc("/api/assistant", assistantHandler)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("IRIS Butler backend running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
