package main

import "testing"

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input      string
		wantAction string
	}{
		{"Hey Butler, change the page to baby pink", "SET_THEME"},
		{"Make this easier to read", "SET_READING_MODE"},
		{"Increase the text size", "SET_TEXT_SIZE"},
		{"Make the page simpler", "SIMPLIFY"},
		{"Turn on high contrast", "HIGH_CONTRAST"},
		{"Reduce motion", "REDUCE_MOTION"},
		{"Reset to default", "RESET"},
		{"Hello butler", ""},
	}

	for _, tt := range tests {
		_, action := parseCommand(tt.input)
		got := ""
		if action != nil {
			got = action.Type
		}
		if got != tt.wantAction {
			t.Errorf("parseCommand(%q) action = %q, want %q", tt.input, got, tt.wantAction)
		}
	}
}
