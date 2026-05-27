package repository

import (
	"testing"
)

func TestRepo_HealthPing(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{
			name:     "Ping OK",
			expected: true,
		},
	}

	repo := NewHealth()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.HealthPing()

			if result != tt.expected {
				t.Errorf("Ping() = %v, want %v", result, tt.expected)
			}
		})
	}
}
