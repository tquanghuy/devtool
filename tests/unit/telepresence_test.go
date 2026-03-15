package unit

import (
	"testing"
	"devtool/internal/checker"
)

func TestParseTelepresenceStatus(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected bool
	}{
		{
			name:     "Daemon running",
			output:   "Root Daemon: Running\nUser Daemon: Running",
			expected: true,
		},
		{
			name:     "Daemon not running",
			output:   "Root Daemon: Not running\nUser Daemon: Not running",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checker.ParseTelepresenceStatus(tc.output)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
