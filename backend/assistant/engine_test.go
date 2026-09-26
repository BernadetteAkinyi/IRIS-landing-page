package assistant

import (
	"testing"

	"backend/models"
	"backend/store"
)

func TestButlerEngine(t *testing.T) {
	st := store.New()
	eng := NewEngine(st)

	tests := []struct {
		name          string
		req           models.AssistantRequest
		expectedInt   string
		expectedAct   string
		expectTheme   string
		expectPayload string
	}{
		{
			name:        "Change page to baby pink (from README)",
			req:         models.AssistantRequest{Transcript: "Hey Butler, change the page to baby pink."},
			expectedInt: "SET_COLOR_THEME",
			expectedAct: "SET_COLOR_THEME",
			expectTheme: "baby-pink",
		},
		{
			name:        "Change page to lavender",
			req:         models.AssistantRequest{Transcript: "Change the page to lavender"},
			expectedInt: "SET_COLOR_THEME",
			expectedAct: "SET_COLOR_THEME",
			expectTheme: "lavender",
		},
		{
			name:        "Make this easier to read (from README)",
			req:         models.AssistantRequest{Transcript: "Hey Butler, make this easier for me to read."},
			expectedInt: "ENABLE_READING_MODE",
			expectedAct: "SET_READING_MODE",
		},
		{
			name: "Make this easier to read with user profile (dyslexia)",
			req: models.AssistantRequest{
				Transcript: "make the text easier to read",
				UserID:     "user_dyslexia",
			},
			expectedInt: "ENABLE_READING_MODE",
			expectedAct: "SET_READING_MODE",
		},
		{
			name:        "Increase text size (from README)",
			req:         models.AssistantRequest{Transcript: "Increase the text size."},
			expectedInt: "INCREASE_TEXT_SIZE",
			expectedAct: "SET_TEXT_SIZE",
		},
		{
			name:        "Turn on text to speech (from README)",
			req:         models.AssistantRequest{Transcript: "Turn on text-to-speech."},
			expectedInt: "START_TTS",
			expectedAct: "SET_TTS",
		},
		{
			name:        "Make the page simpler (from README)",
			req:         models.AssistantRequest{Transcript: "Make the page simpler."},
			expectedInt: "SIMPLIFY_INTERFACE",
			expectedAct: "SET_SIMPLIFIED_LAYOUT",
		},
		{
			name:        "Enable high contrast",
			req:         models.AssistantRequest{Transcript: "Enable high contrast please."},
			expectedInt: "ENABLE_HIGH_CONTRAST",
			expectedAct: "SET_HIGH_CONTRAST",
		},
		{
			name:        "Reduce motion",
			req:         models.AssistantRequest{Transcript: "Reduce motion"},
			expectedInt: "ENABLE_REDUCE_MOTION",
			expectedAct: "SET_REDUCE_MOTION",
		},
		{
			name: "Apply saved preferences with user",
			req: models.AssistantRequest{
				Transcript: "Apply my preferences",
				UserID:     "user_dyslexia",
			},
			expectedInt: "APPLY_USER_PREFERENCES",
			expectedAct: "APPLY_USER_PREFERENCES",
		},
		{
			name:        "Apply saved preferences guest fallback",
			req:         models.AssistantRequest{Transcript: "Apply my preferences"},
			expectedInt: "APPLY_USER_PREFERENCES",
			expectedAct: "APPLY_USER_PREFERENCES",
		},
		{
			name:        "What can you do query",
			req:         models.AssistantRequest{Transcript: "What can you do?"},
			expectedInt: "HELP",
		},
		{
			name:        "Services inquiry",
			req:         models.AssistantRequest{Transcript: "Tell me about healthcare services"},
			expectedInt: "SERVICES_INFO",
		},
		{
			name:        "Reset page command",
			req:         models.AssistantRequest{Transcript: "Reset page"},
			expectedInt: "RESET_INTERFACE",
			expectedAct: "RESET_INTERFACE",
		},
		{
			name:        "Navigate to features",
			req:         models.AssistantRequest{Transcript: "Go to features"},
			expectedInt: "NAVIGATE",
			expectedAct: "NAVIGATE",
		},
		{
			name:        "Unknown query fallback",
			req:         models.AssistantRequest{Transcript: "Can you fly a rocket to Mars?"},
			expectedInt: "UNKNOWN",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := eng.Process(tc.req)
			if resp.Intent != tc.expectedInt {
				t.Errorf("Expected intent %q, got %q (Reply: %s)", tc.expectedInt, resp.Intent, resp.Reply)
			}
			if tc.expectedAct != "" {
				if len(resp.Actions) == 0 {
					t.Fatalf("Expected at least one action of type %q, got none", tc.expectedAct)
				}
				if resp.Actions[0].Type != tc.expectedAct {
					t.Errorf("Expected action type %q, got %q", tc.expectedAct, resp.Actions[0].Type)
				}
			}
			if tc.expectTheme != "" {
				themeVal, ok := resp.Actions[0].Payload["theme"].(string)
				if !ok || themeVal != tc.expectTheme {
					t.Errorf("Expected theme %q, got %v", tc.expectTheme, resp.Actions[0].Payload["theme"])
				}
			}
		})
	}
}
