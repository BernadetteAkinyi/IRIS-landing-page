package models

// AccessibilityPreferences defines all configurable accessibility settings for a user.
type AccessibilityPreferences struct {
	ColorTheme       string `json:"color_theme"`       // e.g. "baby-pink", "lavender", "mint-green", "peach", "sky-blue", "high-contrast", "default"
	ReadingMode      bool   `json:"reading_mode"`      // dyslexia-friendly focus mode
	ReadingFont      string `json:"reading_font"`      // e.g. "Lexend", "OpenDyslexic", "sans-serif"
	TextSize         string `json:"text_size"`         // "small", "medium", "large", "x-large"
	HighContrast     bool   `json:"high_contrast"`     // dark high-contrast mode
	TextToSpeech     bool   `json:"text_to_speech"`    // whether TTS should speak content
	ReduceMotion     bool   `json:"reduce_motion"`     // minimize animations and motion
	SimplifiedLayout bool   `json:"simplified_layout"` // stripped down distraction-free UI
	LineSpacing      string `json:"line_spacing"`      // "normal", "relaxed", "loose"
}

// User represents an IRIS user account.
type User struct {
	ID                 string                   `json:"id"`
	Email              string                   `json:"email"`
	Name               string                   `json:"name"`
	AccessibilityNeeds []string                 `json:"accessibility_needs"` // e.g. ["dyslexia", "low-vision", "adhd", "sensory-sensitivity"]
	Preferences        AccessibilityPreferences `json:"preferences"`
	SavedEvents        []string                 `json:"saved_events"`
}

// ServiceItem represents a health or community service available in IRIS.
type ServiceItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"` // "health", "community", "support"
	Description string `json:"description"`
	Contact     string `json:"contact"`
	Available   bool   `json:"available"`
}

// EventItem represents an accessible community or healthcare event.
type EventItem struct {
	ID                    string   `json:"id"`
	Title                 string   `json:"title"`
	Date                  string   `json:"date"`
	Location              string   `json:"location"`
	Description           string   `json:"description"`
	AccessibilityFeatures []string `json:"accessibility_features"`
}

// ButlerAction instructs the front-end on what visual or functional changes to apply.
type ButlerAction struct {
	Type    string                 `json:"type"`    // e.g. "SET_COLOR_THEME", "SET_READING_MODE", "SET_TEXT_SIZE", etc.
	Payload map[string]interface{} `json:"payload"` // parameters for the action
}

// AssistantRequest is sent from the browser to the Butler assistant.
type AssistantRequest struct {
	Transcript         string                    `json:"transcript"`
	UserID             string                    `json:"user_id,omitempty"`
	Context            string                    `json:"context,omitempty"`
	CurrentPreferences *AccessibilityPreferences `json:"current_preferences,omitempty"`
}

// AssistantResponse is returned to the browser by the Butler assistant.
type AssistantResponse struct {
	Reply       string                    `json:"reply"`
	Actions     []ButlerAction            `json:"actions,omitempty"`
	Preferences *AccessibilityPreferences `json:"preferences,omitempty"`
	Intent      string                    `json:"intent,omitempty"`
}
