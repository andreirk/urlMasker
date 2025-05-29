package main

import (
	"testing"
)

func TestMask(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com", "http://***********"},
		{"visit http://test.com now", "visit http://******** now"},
		{"no links here", "no links here"},
		{"here are two links http://a and http://b", "here are two links http://* and http://*"},
	}

	for _, tt := range tests {
		got := mask(tt.input)
		if got != tt.expected {
			t.Errorf("mask(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
