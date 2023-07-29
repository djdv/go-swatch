package v1

import (
	"math"
	"testing"
	"time"
)

func equalWithTolerance(a, b float64) bool {
	tolerance := 0.000001 // Max length of Swatch Time
	diff := math.Abs(a - b)
	return diff < tolerance
}

func TestNewFromTime(t *testing.T) {
	var now time.Time
	newT := NewFromTimeUsing(now, defaultAlgorithm)
	if newT == nil {
		t.Errorf("expected NewFromTime to return InternetTime")
		return
	}

	if now.UnixNano() != newT.UnixNano() {
		t.Errorf("expected UnixNano of InternetTime and time.Time to match")
	}
}

func TestNew(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected uninitialised InternetTime not to throw deref")
		}
	}()

	newT := NewUsing(defaultAlgorithm)
	if newT == nil {
		t.Errorf("expected New to return InternetTime")
		return
	}

	// Both of these functionally equivalent - just a sanity check
	_ = newT.UnixNano()
	_ = (&internetTime{}).UnixNano()
}

func TestTotalSecondsPreciseBeats(t *testing.T) {
	tests := []struct {
		name   string
		t      string
		expect float64
	}{
		{
			name:   "",
			t:      "2006-02-15T12:00:00-06:00",
			expect: 791.666667,
		},
		{
			name:   "",
			t:      "2008-05-11T03:10:07+10:00",
			expect: 757.025463,
		},
		{
			name:   "",
			t:      "2023-01-02T11:11:28+10:00",
			expect: 91.296296,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := NewFromTime(tTime)
			if beats := newT.PreciseBeats(); !equalWithTolerance(beats, tt.expect) {
				t.Errorf("expect %s to be @%f not @%f",
					tt.t,
					beats,
					tt.expect,
				)
			}
		})
	}
}

func TestTotalNanosecondsPreciseBeats(t *testing.T) {
	tests := []struct {
		name   string
		t      string
		expect float64
	}{
		{
			name:   "",
			t:      "2006-02-15T12:00:00-06:00",
			expect: 791.666667,
		},
		{
			name:   "",
			t:      "2008-05-11T03:10:07+10:00",
			expect: 757.025463,
		},
		{
			name:   "",
			t:      "2023-01-02T11:11:28+10:00",
			expect: 91.296296,
		},
		{
			name:   "",
			t:      "2023-01-02T23:59:59.999999999+02:00",
			expect: 91.296296,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := NewFromTimeUsing(tTime, TotalNanoSeconds)
			if beats := newT.PreciseBeats(); !equalWithTolerance(beats, tt.expect) {
				t.Errorf("expect %s to be @%f not @%f",
					tt.t,
					beats,
					tt.expect,
				)
			}
		})
	}
}
