package unit

import (
	"testing"
	"devtool/internal/checker"
)

func TestParseDockerInfo(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected bool
	}{
		{
			name:     "Daemon running",
			output:   "Client: \n Context: default\nServer: \n Containers: 10\n",
			expected: true,
		},
		{
			name:     "Daemon not running",
			output:   "Client: \n Context: default\nCannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checker.ParseDockerInfo(tc.output)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
