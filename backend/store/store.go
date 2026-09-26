package store

import (
	"errors"
	"strings"
	"sync"

	"backend/models"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
)

type Store struct {
	mu          sync.RWMutex
	users       map[string]*models.User
	usersByMail map[string]*models.User
	services    []models.ServiceItem
	events      map[string]models.EventItem
}

func New() *Store {
	s := &Store{
		users:       make(map[string]*models.User),
		usersByMail: make(map[string]*models.User),
		events:      make(map[string]models.EventItem),
	}
	s.seedData()
	return s
}

func (s *Store) seedData() {
	// Seed demo users with realistic accessibility profiles
	dyslexiaUser := &models.User{
		ID:                 "user_dyslexia",
		Email:              "bernadette@iris.local",
		Name:               "Bernadette",
		AccessibilityNeeds: []string{"dyslexia", "visual-stress"},
		Preferences: models.AccessibilityPreferences{
			ColorTheme:       "baby-pink",
			ReadingMode:      true,
			ReadingFont:      "Lexend",
			TextSize:         "large",
			HighContrast:     false,
			TextToSpeech:     false,
			ReduceMotion:     true,
			SimplifiedLayout: true,
			LineSpacing:      "relaxed",
		},
		SavedEvents: []string{"event_1"},
	}
	s.users[dyslexiaUser.ID] = dyslexiaUser
	s.usersByMail[strings.ToLower(dyslexiaUser.Email)] = dyslexiaUser

	lowVisionUser := &models.User{
		ID:                 "user_low_vision",
		Email:              "alex@iris.local",
		Name:               "Alex",
		AccessibilityNeeds: []string{"low-vision"},
		Preferences: models.AccessibilityPreferences{
			ColorTheme:       "high-contrast",
			ReadingMode:      false,
			ReadingFont:      "Lexend",
			TextSize:         "x-large",
			HighContrast:     true,
			TextToSpeech:     true,
			ReduceMotion:     false,
			SimplifiedLayout: true,
			LineSpacing:      "loose",
		},
		SavedEvents: []string{"event_2"},
	}
	s.users[lowVisionUser.ID] = lowVisionUser
	s.usersByMail[strings.ToLower(lowVisionUser.Email)] = lowVisionUser

	// Seed health and community services
	s.services = []models.ServiceItem{
		{
			ID:          "svc_1",
			Name:        "Accessible Sensory Care Clinic",
			Category:    "health",
			Description: "Low-stimulus sensory health checks and consultations tailored for neurodivergent individuals.",
			Contact:     "sensorycare@iris.local | +1-800-555-0199",
			Available:   true,
		},
		{
			ID:          "svc_2",
			Name:        "Dyslexia & Neurodiversity Digital Resource Hub",
			Category:    "community",
			Description: "Free digital literacy tools, audiobooks, and reading assistance guides.",
			Contact:     "support@iris.local | community.iris.local",
			Available:   true,
		},
		{
			ID:          "svc_3",
			Name:        "Mental Wellbeing & Calm Space Support",
			Category:    "health",
			Description: "Quiet room consultations, mindfulness programs, and anxiety management for invisible disabilities.",
			Contact:     "wellness@iris.local",
			Available:   true,
		},
		{
			ID:          "svc_4",
			Name:        "Inclusive Tech Peer Group",
			Category:    "community",
			Description: "Weekly peer meetings to share adaptive tech hacks and assistive device reviews.",
			Contact:     "groups@iris.local",
			Available:   true,
		},
	}

	// Seed community & health events
	e1 := models.EventItem{
		ID:          "event_1",
		Title:       "Adaptive Digital Accessibility Workshop",
		Date:        "2026-10-15 14:00",
		Location:    "IRIS Virtual Community Center",
		Description: "A hands-on workshop exploring personalized assistive settings, screen reading, and dyslexic-friendly interfaces.",
		AccessibilityFeatures: []string{
			"Live closed captions",
			"Sign language interpreter",
			"Screen-reader verified slides",
			"Sensory break rooms",
		},
	}
	e2 := models.EventItem{
		ID:          "event_2",
		Title:       "Neurodiversity in Everyday Life Meetup",
		Date:        "2026-10-22 17:30",
		Location:    "Greenwood Community Library & Online",
		Description: "Community panel discussing invisible disability navigation and peer support networks.",
		AccessibilityFeatures: []string{
			"Wheelchair accessible",
			"Quiet zone available",
			"High-contrast presentation materials",
		},
	}
	s.events[e1.ID] = e1
	s.events[e2.ID] = e2
}

func (s *Store) GetUser(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	// return copy
	clone := *u
	return &clone, nil
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.usersByMail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return nil, ErrUserNotFound
	}
	clone := *u
	return &clone, nil
}

func (s *Store) CreateOrUpdateUser(u *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if u.ID == "" {
		u.ID = "user_" + strings.ReplaceAll(strings.ToLower(u.Email), "@", "_at_")
	}
	s.users[u.ID] = u
	s.usersByMail[strings.ToLower(u.Email)] = u
	return nil
}

func (s *Store) GetPreferences(userID string) (*models.AccessibilityPreferences, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}
	clone := u.Preferences
	return &clone, nil
}

func (s *Store) UpdatePreferences(userID string, prefs models.AccessibilityPreferences) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[userID]
	if !ok {
		return ErrUserNotFound
	}
	u.Preferences = prefs
	return nil
}

func (s *Store) GetServices(category string) []models.ServiceItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if category == "" {
		res := make([]models.ServiceItem, len(s.services))
		copy(res, s.services)
		return res
	}
	var res []models.ServiceItem
	for _, item := range s.services {
		if strings.EqualFold(item.Category, category) {
			res = append(res, item)
		}
	}
	return res
}

func (s *Store) GetEvents() []models.EventItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]models.EventItem, 0, len(s.events))
	for _, e := range s.events {
		res = append(res, e)
	}
	return res
}

func (s *Store) AddUserEvent(userID, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[userID]
	if !ok {
		return ErrUserNotFound
	}
	if _, ok := s.events[eventID]; !ok {
		return ErrEventNotFound
	}
	for _, id := range u.SavedEvents {
		if id == eventID {
			return nil
		}
	}
	u.SavedEvents = append(u.SavedEvents, eventID)
	return nil
}
