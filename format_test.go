package swatch_test

import (
	"testing"
	"time"

	"github.com/djdv/go-swatch"
)

func TestFormatReturnsBeats(t *testing.T) {
	tests := []struct {
		name          string
		format        string
		expectedValue string
		t             string
	}{
		{
			name:          "Swatch time format",
			format:        swatch.Swatch.String(),
			expectedValue: "@91",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Deci time format",
			format:        swatch.Deci.String(),
			expectedValue: "@91.2",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Centi time format",
			format:        swatch.Centi.String(),
			expectedValue: "@91.29",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Mili time format",
			format:        swatch.Mili.String(),
			expectedValue: "@91.296",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Micro time format",
			format:        swatch.Micro.String(),
			expectedValue: "@91.296296",
			t:             "2023-01-02T11:11:28+10:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := swatch.New(swatch.WithTime(tTime))

			if i := newT.Format(tt.format); i != tt.expectedValue {
				t.Errorf("expected %s got %s", tt.expectedValue, i)
			}
		})
	}
}

func TestInternetTimeString(t *testing.T) {
	tTime, err := time.Parse(time.RFC3339, "2023-01-02T11:11:28+10:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	newT := swatch.New(swatch.WithTime(tTime))

	if s := newT.String(); s != "@91" {
		t.Errorf("output of InternetTime String() unexpected: %s", s)
	}
}

func TestCanCombineFormatting(t *testing.T) {
	tTime, err := time.Parse(time.RFC3339, "2023-01-02T11:11:28+10:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	s := swatch.New(swatch.WithTime(tTime))
	if f := s.Format("2006-01-02 " + swatch.Swatch.String()); f != "2023-01-02 @91" {
		t.Errorf("Failed to mix formating, got %s", f)
	}
}
