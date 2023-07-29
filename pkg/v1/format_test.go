package v1

import (
	"testing"
	"time"
)

func TestFormatReturnsBeats(t *testing.T) {
	tests := []struct {
		name          string
		format        string
		expectedValue string
		t             time.Time
	}{
		{
			name:          "Swatch time format",
			format:        Swatch,
			expectedValue: "",
			t:             time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iTime := NewFromTime(tt.t)
			if i := iTime.Format(tt.format); i != tt.expectedValue {
				t.Errorf("expected %s got %s", tt.expectedValue, i)
			}
		})
	}
}
