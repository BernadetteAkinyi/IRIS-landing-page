package assistant

import (
	"fmt"
	"regexp"
	"strings"

	"backend/models"
	"backend/store"
)

type Engine struct {
	store *store.Store
}

func NewEngine(s *store.Store) *Engine {
	return &Engine{store: s}
}

// ColorThemeDefinition maps friendly color names to CSS classes and hex codes
type ColorThemeDefinition struct {
	ID        string
	Name      string
	Hex       string
	TextColor string
}

var supportedThemes = map[string]ColorThemeDefinition{
	"baby-pink":   {ID: "baby-pink", Name: "Baby Pink", Hex: "#FAD2E1", TextColor: "#4A2234"},
	"lavender":    {ID: "lavender", Name: "Lavender", Hex: "#E2D9F3", TextColor: "#33224B"},
	"mint-green":  {ID: "mint-green", Name: "Mint Green", Hex: "#D8F3DC", TextColor: "#1B4332"},
	"peach":       {ID: "peach", Name: "Soft Peach", Hex: "#FFE5D9", TextColor: "#582F24"},
	"sky-blue":    {ID: "sky-blue", Name: "Sky Blue", Hex: "#D0E8FF", TextColor: "#1E3A5F"},
	"high-contrast": {ID: "high-contrast", Name: "High Contrast Dark", Hex: "#121212", TextColor: "#FFFFFF"},
	"default":     {ID: "default", Name: "Default Warm Sand", Hex: "#FDF4EE", TextColor: "#4A2E2A"},
}

// Process analyzes the incoming request transcript and generates an assistant response with actionable commands.
func (e *Engine) Process(req models.AssistantRequest) models.AssistantResponse {
	transcript := strings.TrimSpace(req.Transcript)
	tLower := strings.ToLower(transcript)

	// Clean punctuation into spaces (replacing hyphens, commas, periods etc.)
	clean := regexp.MustCompile(`[^\w\s]`).ReplaceAllString(tLower, " ")
	clean = regexp.MustCompile(`\s+`).ReplaceAllString(clean, " ")
	clean = strings.TrimSpace(clean)

	var user *models.User
	if req.UserID != "" && e.store != nil {
		user, _ = e.store.GetUser(req.UserID)
	}

	// 1. Color theme commands: "change the page to baby pink", "make it lavender", "set color to mint"
	if themeResp, matched := e.matchColorTheme(clean); matched {
		return themeResp
	}

	// 2. Reading mode commands: "make this easier to read", "make the text easier to read", "dyslexia mode", "focus reading"
	if readingResp, matched := e.matchReadingMode(clean, user); matched {
		return readingResp
	}

	// 3. Text size commands: "increase the text size", "make text bigger", "larger font", "smaller text", "reset text size"
	if textSizeResp, matched := e.matchTextSize(clean, req.CurrentPreferences); matched {
		return textSizeResp
	}

	// 4. Text-to-speech commands: "turn on text-to-speech", "read to me", "read the page", "stop reading", "mute"
	if ttsResp, matched := e.matchTextToSpeech(clean); matched {
		return ttsResp
	}

	// 5. Simplify interface commands: "make the page simpler", "simplify the interface", "clean layout", "standard view"
	if simplifyResp, matched := e.matchSimplify(clean); matched {
		return simplifyResp
	}

	// 6. High contrast commands: "enable high contrast", "turn on high contrast", "high contrast mode", "dark mode"
	if contrastResp, matched := e.matchHighContrast(clean); matched {
		return contrastResp
	}

	// 7. Reduce motion commands: "reduce motion", "stop animation", "calm mode", "enable motion"
	if motionResp, matched := e.matchReduceMotion(clean); matched {
		return motionResp
	}

	// 8. Apply saved user profile / preferences
	if profileResp, matched := e.matchApplyProfile(clean, user, req.UserID); matched {
		return profileResp
	}

	// 9. Reset interface commands: "reset page", "default view", "clear settings"
	if resetResp, matched := e.matchReset(clean); matched {
		return resetResp
	}

	// 10. Navigation commands: "go to features", "go to home", "go to login", "contact"
	if navResp, matched := e.matchNavigation(clean); matched {
		return navResp
	}

	// 11. Informational queries: Services, Events, About IRIS, Capabilities, Help, Greetings
	if infoResp, matched := e.matchInformation(clean); matched {
		return infoResp
	}

	// Default polite Butler fallback
	return models.AssistantResponse{
		Reply:  fmt.Sprintf("I heard: \"%s\". While I am not certain how to execute that request, I can adjust the page color, make text easier to read, alter font sizes, read content aloud, simplify the layout, or apply your accessibility profile. How may I be of assistance?", transcript),
		Intent: "UNKNOWN",
	}
}

func (e *Engine) matchColorTheme(clean string) (models.AssistantResponse, bool) {
	// Check for specific themes
	colorPatterns := []struct {
		keywords []string
		themeKey string
	}{
		{[]string{"baby pink", "pink", "rose", "blush"}, "baby-pink"},
		{[]string{"lavender", "purple", "violet", "lilac"}, "lavender"},
		{[]string{"mint green", "mint", "pale green", "green"}, "mint-green"},
		{[]string{"peach", "soft peach", "apricot", "warm coral"}, "peach"},
		{[]string{"sky blue", "light blue", "cyan", "baby blue", "blue"}, "sky-blue"},
		{[]string{"high contrast", "dark mode", "black theme"}, "high-contrast"},
		{[]string{"default color", "default theme", "sand", "warm sand", "original color"}, "default"},
	}

	isColorQuery := strings.Contains(clean, "color") ||
		strings.Contains(clean, "theme") ||
		strings.Contains(clean, "background") ||
		strings.Contains(clean, "change the page to") ||
		strings.Contains(clean, "turn the page to") ||
		strings.Contains(clean, "make the page") ||
		strings.Contains(clean, "switch to") ||
		strings.Contains(clean, "palette")

	for _, cp := range colorPatterns {
		for _, kw := range cp.keywords {
			if strings.Contains(clean, kw) && (isColorQuery || strings.HasPrefix(clean, "make it "+kw) || clean == kw) {
				theme := supportedThemes[cp.themeKey]
				return models.AssistantResponse{
					Reply: fmt.Sprintf("Right away. I have changed the page color to %s.", theme.Name),
					Actions: []models.ButlerAction{
						{
							Type: "SET_COLOR_THEME",
							Payload: map[string]interface{}{
								"theme":      theme.ID,
								"name":       theme.Name,
								"hex":        theme.Hex,
								"text_color": theme.TextColor,
							},
						},
					},
					Preferences: &models.AccessibilityPreferences{
						ColorTheme: theme.ID,
					},
					Intent: "SET_COLOR_THEME",
				}, true
			}
		}
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchReadingMode(clean string, user *models.User) (models.AssistantResponse, bool) {
	triggers := []string{
		"make this easier to read",
		"make this easier for me to read",
		"make it easier for me to read",
		"make it easier to read",
		"make the text easier to read",
		"make text easier to read",
		"easier for me to read",
		"easier to read",
		"reading mode",
		"dyslexia",
		"dyslexic",
		"focus reading",
		"reading preferences",
		"improve readability",
		"clean typography",
	}

	disableTriggers := []string{
		"disable reading mode",
		"turn off reading mode",
		"exit reading mode",
		"stop reading mode",
	}

	for _, dt := range disableTriggers {
		if strings.Contains(clean, dt) {
			return models.AssistantResponse{
				Reply: "Understood. I have turned off reading mode for you.",
				Actions: []models.ButlerAction{
					{
						Type: "SET_READING_MODE",
						Payload: map[string]interface{}{
							"enabled": false,
						},
					},
				},
				Preferences: &models.AccessibilityPreferences{
					ReadingMode: false,
				},
				Intent: "DISABLE_READING_MODE",
			}, true
		}
	}

	for _, t := range triggers {
		if strings.Contains(clean, t) {
			// If user has saved preferences, customize according to their profile
			font := "Lexend"
			spacing := "relaxed"
			reply := "Certainly. I have activated reading mode with enhanced letter spacing and clear typography to make the text easier to read."

			if user != nil && user.Preferences.ReadingFont != "" {
				font = user.Preferences.ReadingFont
				if user.Preferences.LineSpacing != "" {
					spacing = user.Preferences.LineSpacing
				}
				reply = fmt.Sprintf("Certainly. I have applied your saved reading preferences using %s typography and %s spacing.", font, spacing)
			}

			return models.AssistantResponse{
				Reply: reply,
				Actions: []models.ButlerAction{
					{
						Type: "SET_READING_MODE",
						Payload: map[string]interface{}{
							"enabled":      true,
							"font":         font,
							"line_spacing": spacing,
						},
					},
				},
				Preferences: &models.AccessibilityPreferences{
					ReadingMode:  true,
					ReadingFont:  font,
					LineSpacing:  spacing,
				},
				Intent: "ENABLE_READING_MODE",
			}, true
		}
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchTextSize(clean string, current *models.AccessibilityPreferences) (models.AssistantResponse, bool) {
	// Increase text size
	if strings.Contains(clean, "increase text") ||
		strings.Contains(clean, "increase the text") ||
		strings.Contains(clean, "bigger text") ||
		strings.Contains(clean, "larger text") ||
		strings.Contains(clean, "larger font") ||
		strings.Contains(clean, "bigger font") ||
		strings.Contains(clean, "make text larger") ||
		strings.Contains(clean, "make text bigger") ||
		strings.Contains(clean, "make font bigger") {

		targetSize := "large"
		scale := 1.2
		if current != nil && current.TextSize == "large" {
			targetSize = "x-large"
			scale = 1.35
		}

		return models.AssistantResponse{
			Reply: fmt.Sprintf("At once. I have increased the text size to %s.", targetSize),
			Actions: []models.ButlerAction{
				{
					Type: "SET_TEXT_SIZE",
					Payload: map[string]interface{}{
						"size":  targetSize,
						"scale": scale,
					},
				},
			},
			Preferences: &models.AccessibilityPreferences{
				TextSize: targetSize,
			},
			Intent: "INCREASE_TEXT_SIZE",
		}, true
	}

	// Decrease text size
	if strings.Contains(clean, "decrease text") ||
		strings.Contains(clean, "decrease the text") ||
		strings.Contains(clean, "smaller text") ||
		strings.Contains(clean, "smaller font") ||
		strings.Contains(clean, "make text smaller") ||
		strings.Contains(clean, "make font smaller") ||
		strings.Contains(clean, "reduce text size") {

		targetSize := "small"
		scale := 0.9
		return models.AssistantResponse{
			Reply: "Certainly. I have reduced the text size.",
			Actions: []models.ButlerAction{
				{
					Type: "SET_TEXT_SIZE",
					Payload: map[string]interface{}{
						"size":  targetSize,
						"scale": scale,
					},
				},
			},
			Preferences: &models.AccessibilityPreferences{
				TextSize: targetSize,
			},
			Intent: "DECREASE_TEXT_SIZE",
		}, true
	}

	// Reset text size
	if strings.Contains(clean, "reset text size") ||
		strings.Contains(clean, "normal text size") ||
		strings.Contains(clean, "default text size") ||
		strings.Contains(clean, "standard font size") {
		return models.AssistantResponse{
			Reply: "I have restored the text size to standard.",
			Actions: []models.ButlerAction{
				{
					Type: "SET_TEXT_SIZE",
					Payload: map[string]interface{}{
						"size":  "medium",
						"scale": 1.0,
					},
				},
			},
			Preferences: &models.AccessibilityPreferences{
				TextSize: "medium",
			},
			Intent: "RESET_TEXT_SIZE",
		}, true
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchTextToSpeech(clean string) (models.AssistantResponse, bool) {
	stopKeywords := []string{
		"stop speaking",
		"stop talking",
		"stop reading",
		"mute",
		"turn off text to speech",
		"disable text to speech",
		"silence",
		"quiet",
	}
	for _, kw := range stopKeywords {
		if strings.Contains(clean, kw) {
			return models.AssistantResponse{
				Reply: "Very well, I shall remain quiet.",
				Actions: []models.ButlerAction{
					{
						Type: "SET_TTS",
						Payload: map[string]interface{}{
							"enabled": false,
							"action":  "stop",
						},
					},
				},
				Preferences: &models.AccessibilityPreferences{
					TextToSpeech: false,
				},
				Intent: "STOP_TTS",
			}, true
		}
	}

	startKeywords := []string{
		"turn on text to speech",
		"enable text to speech",
		"read this to me",
		"read to me",
		"read the page",
		"read out loud",
		"read aloud",
		"speak the content",
		"speak",
	}
	for _, kw := range startKeywords {
		if strings.Contains(clean, kw) {
			return models.AssistantResponse{
				Reply: "Of course. Text-to-speech is activated. I will read the on-screen content for you.",
				Actions: []models.ButlerAction{
					{
						Type: "SET_TTS",
						Payload: map[string]interface{}{
							"enabled":      true,
							"action":       "read",
							"read_current": true,
						},
					},
				},
				Preferences: &models.AccessibilityPreferences{
					TextToSpeech: true,
				},
				Intent: "START_TTS",
			}, true
		}
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchSimplify(clean string) (models.AssistantResponse, bool) {
	if strings.Contains(clean, "make the page simpler") ||
		strings.Contains(clean, "make page simpler") ||
		strings.Contains(clean, "simplify the interface") ||
		strings.Contains(clean, "simplify the page") ||
		strings.Contains(clean, "simplify") ||
		strings.Contains(clean, "minimal view") ||
		strings.Contains(clean, "clean view") ||
		strings.Contains(clean, "distraction free") {
		return models.AssistantResponse{
			Reply: "Certainly. I have simplified the interface to remove visual clutter and highlight essential content.",
			Actions: []models.ButlerAction{
				{
					Type: "SET_SIMPLIFIED_LAYOUT",
					Payload: map[string]interface{}{
						"enabled": true,
					},
				},
			},
			Preferences: &models.AccessibilityPreferences{
				SimplifiedLayout: true,
			},
			Intent: "SIMPLIFY_INTERFACE",
		}, true
	}

	if strings.Contains(clean, "restore layout") ||
		strings.Contains(clean, "standard layout") ||
		strings.Contains(clean, "full layout") ||
		strings.Contains(clean, "turn off simplify") {
		return models.AssistantResponse{
			Reply: "I have restored the full standard layout.",
			Actions: []models.ButlerAction{
				{
					Type: "SET_SIMPLIFIED_LAYOUT",
					Payload: map[string]interface{}{
						"enabled": false,
					},
				},
			},
			Preferences: &models.AccessibilityPreferences{
				SimplifiedLayout: false,
			},
			Intent: "RESTORE_LAYOUT",
		}, true
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchHighContrast(clean string) (models.AssistantResponse, bool) {
	disableTriggers := []string{
		"turn off high contrast",
		"disable high contrast",
		"normal contrast",
		"standard contrast",
	}
	for _, dt := range disableTriggers {
		if strings.Contains(clean, dt) {
			return models.AssistantResponse{
				Reply: "High contrast mode has been turned off.",
				Actions: []models.ButlerAction{
					{
						Type: "SET_HIGH_CONTRAST",
						Payload: map[string]interface{}{
							"enabled": false,
						},
					},
				},
				Preferences: &models.AccessibilityPreferences{
					HighContrast: false,
				},
				Intent: "DISABLE_HIGH_CONTRAST",
			}, true
		}
	}

	enableTriggers := []string{
		"high contrast",
		"enable high contrast",
		"turn on high contrast",
		"contrast mode",
		"dark contrast",
	}
	for _, et := range enableTriggers {
		if strings.Contains(clean, et) {
			return models.AssistantResponse{
				Reply: "High contrast has been enabled for maximum readability and visual definition.",
				Actions: []models.ButlerAction{
					{
						Type: "SET_HIGH_CONTRAST",
						Payload: map[string]interface{}{
							"enabled": true,
						},
					},
				},
				Preferences: &models.AccessibilityPreferences{
					HighContrast: true,
				},
				Intent: "ENABLE_HIGH_CONTRAST",
			}, true
		}
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchReduceMotion(clean string) (models.AssistantResponse, bool) {
	if strings.Contains(clean, "reduce motion") ||
		strings.Contains(clean, "stop animation") ||
		strings.Contains(clean, "stop animations") ||
		strings.Contains(clean, "calmer mode") ||
		strings.Contains(clean, "less motion") {
		return models.AssistantResponse{
			Reply: "Reduced motion is now enabled. Animations and transitions have been minimized.",
			Actions: []models.ButlerAction{
				{
					Type: "SET_REDUCE_MOTION",
					Payload: map[string]interface{}{
						"enabled": true,
					},
				},
			},
			Preferences: &models.AccessibilityPreferences{
				ReduceMotion: true,
			},
			Intent: "ENABLE_REDUCE_MOTION",
		}, true
	}

	if strings.Contains(clean, "enable motion") ||
		strings.Contains(clean, "allow animation") ||
		strings.Contains(clean, "restore animation") {
		return models.AssistantResponse{
			Reply: "Motion animations have been restored.",
			Actions: []models.ButlerAction{
				{
					Type: "SET_REDUCE_MOTION",
					Payload: map[string]interface{}{
						"enabled": false,
					},
				},
			},
			Preferences: &models.AccessibilityPreferences{
				ReduceMotion: false,
			},
			Intent: "DISABLE_REDUCE_MOTION",
		}, true
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchApplyProfile(clean string, user *models.User, userID string) (models.AssistantResponse, bool) {
	triggers := []string{
		"apply my preferences",
		"apply my profile",
		"load my preferences",
		"load my profile",
		"use my settings",
		"my preferences",
		"personalize for me",
	}

	for _, t := range triggers {
		if strings.Contains(clean, t) {
			if user != nil {
				return models.AssistantResponse{
					Reply: fmt.Sprintf("Welcome, %s. I have applied your personalized accessibility profile: %s theme, reading mode %v, and %s text size.",
						user.Name, user.Preferences.ColorTheme, user.Preferences.ReadingMode, user.Preferences.TextSize),
					Actions: []models.ButlerAction{
						{
							Type: "APPLY_USER_PREFERENCES",
							Payload: map[string]interface{}{
								"preferences": user.Preferences,
								"user_name":   user.Name,
							},
						},
					},
					Preferences: &user.Preferences,
					Intent:      "APPLY_USER_PREFERENCES",
				}, true
			}

			// If no user was passed, use default personalized profile
			defaultPrefs := models.AccessibilityPreferences{
				ColorTheme:       "baby-pink",
				ReadingMode:      true,
				ReadingFont:      "Lexend",
				TextSize:         "large",
				HighContrast:     false,
				TextToSpeech:     false,
				ReduceMotion:     true,
				SimplifiedLayout: false,
				LineSpacing:      "relaxed",
			}
			return models.AssistantResponse{
				Reply: "I have applied the default IRIS accessibility profile: gentle baby-pink theme, enhanced Lexend reading typography, and relaxed spacing.",
				Actions: []models.ButlerAction{
					{
						Type: "APPLY_USER_PREFERENCES",
						Payload: map[string]interface{}{
							"preferences": defaultPrefs,
							"user_name":   "Guest",
						},
					},
				},
				Preferences: &defaultPrefs,
				Intent:      "APPLY_USER_PREFERENCES",
			}, true
		}
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchReset(clean string) (models.AssistantResponse, bool) {
	if strings.Contains(clean, "reset page") ||
		strings.Contains(clean, "reset everything") ||
		strings.Contains(clean, "reset interface") ||
		strings.Contains(clean, "default view") ||
		strings.Contains(clean, "clear changes") ||
		strings.Contains(clean, "reset settings") {
		return models.AssistantResponse{
			Reply: "I have restored all page settings to their original defaults.",
			Actions: []models.ButlerAction{
				{
					Type:    "RESET_INTERFACE",
					Payload: map[string]interface{}{},
				},
			},
			Intent: "RESET_INTERFACE",
		}, true
	}
	return models.AssistantResponse{}, false
}

func (e *Engine) matchNavigation(clean string) (models.AssistantResponse, bool) {
	if strings.Contains(clean, "go to features") ||
		strings.Contains(clean, "show features") ||
		strings.Contains(clean, "open features") {
		return models.AssistantResponse{
			Reply: "Navigating you to the Features page right away.",
			Actions: []models.ButlerAction{
				{
					Type: "NAVIGATE",
					Payload: map[string]interface{}{
						"page": "features.html",
					},
				},
			},
			Intent: "NAVIGATE",
		}, true
	}

	if strings.Contains(clean, "go home") ||
		strings.Contains(clean, "go to home") ||
		strings.Contains(clean, "open home") {
		return models.AssistantResponse{
			Reply: "Navigating to the Home page.",
			Actions: []models.ButlerAction{
				{
					Type: "NAVIGATE",
					Payload: map[string]interface{}{
						"page": "index.html",
					},
				},
			},
			Intent: "NAVIGATE",
		}, true
	}

	if strings.Contains(clean, "go to login") ||
		strings.Contains(clean, "open login") ||
		strings.Contains(clean, "sign in") {
		return models.AssistantResponse{
			Reply: "Directing you to the Sign-in page.",
			Actions: []models.ButlerAction{
				{
					Type: "NAVIGATE",
					Payload: map[string]interface{}{
						"page": "login.html",
					},
				},
			},
			Intent: "NAVIGATE",
		}, true
	}

	return models.AssistantResponse{}, false
}

func (e *Engine) matchInformation(clean string) (models.AssistantResponse, bool) {
	// Greetings
	if clean == "hello" || clean == "hi" || clean == "hey" ||
		clean == "hello butler" || clean == "hi butler" || clean == "hey butler" ||
		clean == "good day" ||
		strings.Contains(clean, "good morning") ||
		strings.Contains(clean, "good afternoon") ||
		strings.Contains(clean, "good evening") {
		return models.AssistantResponse{
			Reply:  "Good day to you. I am IRIS Butler, your accessibility companion. How may I be of service?",
			Intent: "GREETING",
		}, true
	}

	// Gratitude
	if strings.Contains(clean, "thank") {
		return models.AssistantResponse{
			Reply:  "It is entirely my pleasure. Never hesitate to ask if there is anything else you require.",
			Intent: "GRATITUDE",
		}, true
	}

	// What can you do / capabilities / help
	if strings.Contains(clean, "what can you do") ||
		strings.Contains(clean, "who are you") ||
		strings.Contains(clean, "help") ||
		strings.Contains(clean, "capabilities") ||
		strings.Contains(clean, "instructions") {
		return models.AssistantResponse{
			Reply: "I am Butler, your adaptive accessibility companion. You can ask me to change the page color (such as 'change the page to baby pink' or 'lavender'), make text easier to read, increase or decrease text size, turn on text-to-speech, enable high contrast, simplify the layout, or apply your saved profile.",
			Intent: "HELP",
		}, true
	}

	// About IRIS
	if strings.Contains(clean, "about iris") || strings.Contains(clean, "what is iris") {
		return models.AssistantResponse{
			Reply: "IRIS is an adaptive accessibility companion designed for individuals with invisible disabilities. Instead of forcing you to adapt to a rigid interface, IRIS tailors colors, typography, layout, and sensory feedback to your individual preferences.",
			Intent: "ABOUT_IRIS",
		}, true
	}

	// Available health / community services
	if strings.Contains(clean, "services") || strings.Contains(clean, "health care") || strings.Contains(clean, "healthcare") {
		var summary string
		if e.store != nil {
			services := e.store.GetServices("")
			summary = fmt.Sprintf("IRIS offers %d integrated services, including the Accessible Sensory Care Clinic and the Dyslexia Digital Resource Hub.", len(services))
		} else {
			summary = "IRIS connects you with accessible health and community services tailored to your needs."
		}
		return models.AssistantResponse{
			Reply:  summary + " Would you like me to guide you to our features or details?",
			Intent: "SERVICES_INFO",
		}, true
	}

	// Upcoming events
	if strings.Contains(clean, "event") || strings.Contains(clean, "events") || strings.Contains(clean, "workshop") {
		var summary string
		if e.store != nil {
			events := e.store.GetEvents()
			summary = fmt.Sprintf("There are %d upcoming accessible events, including the Adaptive Digital Accessibility Workshop.", len(events))
		} else {
			summary = "We host regular inclusive community events and workshops."
		}
		return models.AssistantResponse{
			Reply:  summary,
			Intent: "EVENTS_INFO",
		}, true
	}

	// Time query
	if strings.Contains(clean, "time") {
		return models.AssistantResponse{
			Reply:  "I'm at your service at all hours; every moment is an opportune time to make digital spaces more accessible.",
			Intent: "TIME",
		}, true
	}

	return models.AssistantResponse{}, false
}
