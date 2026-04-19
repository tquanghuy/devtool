package commands

import (
	"testing"
)

func TestFormatCommand(t *testing.T) {
	osCmd := NewOSCommand()

	tests := []struct {
		cmd      string
		port     int
		expected string
	}{
		{"echo %d", 8080, "echo 8080"},
		{"curl localhost:%d/health", 3000, "curl localhost:3000/health"},
		{"docker run -p %d:80 nginx", 80, "docker run -p 80:80 nginx"},
		{"no placeholder", 1234, "no placeholder"},
	}

	for _, tt := range tests {
		result := osCmd.FormatCommand(tt.cmd, tt.port)
		if result != tt.expected {
			t.Errorf("FormatCommand(%q, %d) = %q; want %q", tt.cmd, tt.port, result, tt.expected)
		}
	}
}

func TestGetFreePort(t *testing.T) {
	osCmd := NewOSCommand()

	// Test finding a free port (should return at least the startPort if free)
	port := osCmd.GetFreePort(10000)
	if port < 10000 {
		t.Errorf("GetFreePort(10000) = %d; want >= 10000", port)
	}
}
