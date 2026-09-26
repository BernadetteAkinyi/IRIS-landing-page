package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/assistant"
	"backend/models"
	"backend/store"
)

type API struct {
	store  *store.Store
	engine *assistant.Engine
}

func NewAPI(s *store.Store, e *assistant.Engine) *API {
	return &API{
		store:  s,
		engine: e,
	}
}

// EnableCORS sets headers to allow cross-origin requests from frontend
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// AssistantHandler handles voice and text commands sent to the Butler
func (api *API) AssistantHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AssistantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp := api.engine.Process(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PreferencesHandler retrieves or updates accessibility preferences
func (api *API) PreferencesHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "user_dyslexia" // default demo profile
	}

	switch r.Method {
	case http.MethodGet:
		prefs, err := api.store.GetPreferences(userID)
		if err != nil {
			// Return default preferences if user not found
			defaultPrefs := models.AccessibilityPreferences{
				ColorTheme:       "baby-pink",
				ReadingMode:      true,
				ReadingFont:      "Lexend",
				TextSize:         "medium",
				HighContrast:     false,
				TextToSpeech:     false,
				ReduceMotion:     false,
				SimplifiedLayout: false,
				LineSpacing:      "normal",
			}
			json.NewEncoder(w).Encode(defaultPrefs)
			return
		}
		json.NewEncoder(w).Encode(prefs)

	case http.MethodPut, http.MethodPost:
		var prefs models.AccessibilityPreferences
		if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		_ = api.store.UpdatePreferences(userID, prefs)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "success",
			"message":     "Preferences updated successfully",
			"preferences": prefs,
		})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// LoginHandler simulates user authentication and returns saved profile with preferences
func (api *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	u, err := api.store.GetUserByEmail(creds.Email)
	if err != nil {
		// Create a dynamic profile for the login if not found yet
		name := strings.Split(creds.Email, "@")[0]
		u = &models.User{
			ID:    "user_" + name,
			Email: creds.Email,
			Name:  strings.Title(name),
			AccessibilityNeeds: []string{"adaptive-ui"},
			Preferences: models.AccessibilityPreferences{
				ColorTheme:       "baby-pink",
				ReadingMode:      false,
				ReadingFont:      "Lexend",
				TextSize:         "medium",
				HighContrast:     false,
				TextToSpeech:     false,
				ReduceMotion:     false,
				SimplifiedLayout: false,
				LineSpacing:      "normal",
			},
		}
		_ = api.store.CreateOrUpdateUser(u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":   u,
	})
}

// RegisterHandler allows sign-up with accessibility needs & preferences
func (api *API) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email              string                          `json:"email"`
		Name               string                          `json:"name"`
		AccessibilityNeeds []string                        `json:"accessibility_needs"`
		Preferences        models.AccessibilityPreferences `json:"preferences"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:              req.Email,
		Name:               req.Name,
		AccessibilityNeeds: req.AccessibilityNeeds,
		Preferences:        req.Preferences,
	}
	_ = api.store.CreateOrUpdateUser(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":   user,
	})
}

// ServicesHandler returns health and community services
func (api *API) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	category := r.URL.Query().Get("category")
	services := api.store.GetServices(category)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// EventsHandler returns accessible community and health events
func (api *API) EventsHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	events := api.store.GetEvents()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
