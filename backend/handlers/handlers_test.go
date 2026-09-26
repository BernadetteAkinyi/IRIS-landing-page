package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/assistant"
	"backend/models"
	"backend/store"
)

func setupTestAPI() *API {
	st := store.New()
	eng := assistant.NewEngine(st)
	return NewAPI(st, eng)
}

func TestAssistantHandler(t *testing.T) {
	api := setupTestAPI()

	// Test baby pink command
	payload := models.AssistantRequest{
		Transcript: "Hey Butler, change the page to baby pink.",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/assistant", bytes.NewReader(body))
	w := httptest.NewRecorder()

	api.AssistantHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp models.AssistantResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Intent != "SET_COLOR_THEME" {
		t.Errorf("Expected intent SET_COLOR_THEME, got %s", resp.Intent)
	}
	if len(resp.Actions) == 0 || resp.Actions[0].Type != "SET_COLOR_THEME" {
		t.Errorf("Expected SET_COLOR_THEME action, got %v", resp.Actions)
	}
}

func TestPreferencesHandler(t *testing.T) {
	api := setupTestAPI()

	// Test GET preferences
	req := httptest.NewRequest(http.MethodGet, "/api/preferences?user_id=user_dyslexia", nil)
	w := httptest.NewRecorder()

	api.PreferencesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var prefs models.AccessibilityPreferences
	if err := json.NewDecoder(w.Body).Decode(&prefs); err != nil {
		t.Fatalf("Failed to decode preferences: %v", err)
	}

	if prefs.ColorTheme != "baby-pink" {
		t.Errorf("Expected color theme baby-pink, got %s", prefs.ColorTheme)
	}

	// Test PUT preferences
	prefs.TextSize = "x-large"
	putBody, _ := json.Marshal(prefs)
	putReq := httptest.NewRequest(http.MethodPut, "/api/preferences?user_id=user_dyslexia", bytes.NewReader(putBody))
	putW := httptest.NewRecorder()

	api.PreferencesHandler(putW, putReq)
	if putW.Code != http.StatusOK {
		t.Fatalf("Expected status 200 on PUT, got %d", putW.Code)
	}
}

func TestServicesAndEventsHandler(t *testing.T) {
	api := setupTestAPI()

	// Test GET services
	reqSvc := httptest.NewRequest(http.MethodGet, "/api/services", nil)
	wSvc := httptest.NewRecorder()
	api.ServicesHandler(wSvc, reqSvc)
	if wSvc.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for services, got %d", wSvc.Code)
	}

	var services []models.ServiceItem
	if err := json.NewDecoder(wSvc.Body).Decode(&services); err != nil {
		t.Fatalf("Failed to decode services: %v", err)
	}
	if len(services) == 0 {
		t.Errorf("Expected at least one service, got 0")
	}

	// Test GET events
	reqEvt := httptest.NewRequest(http.MethodGet, "/api/events", nil)
	wEvt := httptest.NewRecorder()
	api.EventsHandler(wEvt, reqEvt)
	if wEvt.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for events, got %d", wEvt.Code)
	}

	var events []models.EventItem
	if err := json.NewDecoder(wEvt.Body).Decode(&events); err != nil {
		t.Fatalf("Failed to decode events: %v", err)
	}
	if len(events) == 0 {
		t.Errorf("Expected at least one event, got 0")
	}
}
