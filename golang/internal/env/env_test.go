package env

import (
	"os"
	"testing"
)

func TestGetOrElse(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "returns env var when set",
			key:          "TEST_VAR",
			envValue:     "test_value",
			defaultValue: "default",
			expected:     "test_value",
		},
		{
			name:         "returns default when env var not set",
			key:          "NONEXISTENT_VAR",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := GetOrElse(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetIntOrElse(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:     "returns parsed int when env var is valid",
			key:      "TEST_VAR_INT",
			envValue: "42",
			expected: 42,
		},
		{
			name:         "returns default when env var not set",
			key:          "TEST_VAR_INT",
			envValue:     "",
			defaultValue: 4242,
			expected:     4242,
		},
		{
			name:         "returns default when env var is invalid",
			key:          "NOT_A_NUMBER",
			envValue:     "0xabcd",
			defaultValue: 4242,
			expected:     4242,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}
			result := GetIntOrElse(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
