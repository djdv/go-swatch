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

func TestSanity(t *testing.T) {
	// When given two dates exactly a beat apart, the beats are indeed 1 beat difference
	t1, err := time.Parse(time.RFC3339, "2006-02-15T12:00:00.000+01:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	t2, err := time.Parse(time.RFC3339, "2006-02-15T12:01:26.4000+01:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	i1 := NewFromTime(t1)
	i2 := NewFromTime(t2)

	a := roundDownFloat(i1.PreciseBeats(), 0)
	b := roundDownFloat(i2.PreciseBeats(), 0)

	// Just to test sanity of GetTime
	if t1.Format("2006-01-02") != i1.GetTime().Format("2006-01-02") {
		t.Errorf("expected t1 date to be the same as swatchTime date")
	}

	if (b - a) < 1 {
		t.Errorf("expected b to be exactly 1 increment higher than a. Got a: %f b: %f", a, b)
	}
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
		{
			name:   "",
			t:      "2023-01-02T23:59:59.999999+01:00",
			expect: 999.988426,
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
					tt.expect,
					beats,
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
			expect: 791.666666,
		},
		{
			name:   "",
			t:      "2008-05-11T03:10:07+10:00",
			expect: 757.025462,
		},
		{
			name:   "",
			t:      "2023-01-02T11:11:28+10:00",
			expect: 91.296296,
		},
		{
			name:   "",
			t:      "2023-01-02T23:59:59.999999999+01:00",
			expect: 999.999999,
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
					tt.expect,
					beats,
				)
			}
		})
	}
}

func TestTotalSecondsBeats(t *testing.T) {
	tests := []struct {
		name   string
		t      string
		expect int
	}{
		{
			name:   "",
			t:      "2006-02-15T12:00:00-06:00",
			expect: 791,
		},
		{
			name:   "",
			t:      "2008-05-11T03:10:07+10:00",
			expect: 757,
		},
		{
			name:   "",
			t:      "2023-01-02T11:11:28+10:00",
			expect: 91,
		},
		{
			name:   "",
			t:      "2023-01-02T23:59:59.999999+01:00",
			expect: 999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := NewFromTime(tTime)
			if beats := newT.Beats(); beats != tt.expect {
				t.Errorf("expect %s to be @%d not @%d",
					tt.t,
					tt.expect,
					beats,
				)
			}
		})
	}
}

func TestTotalNanosecondsBeats(t *testing.T) {
	tests := []struct {
		name   string
		t      string
		expect int
	}{
		{
			name:   "",
			t:      "2006-02-15T12:00:00-06:00",
			expect: 791,
		},
		{
			name:   "",
			t:      "2008-05-11T03:10:07+10:00",
			expect: 757,
		},
		{
			name:   "",
			t:      "2023-01-02T11:11:28+10:00",
			expect: 91,
		},
		{
			name:   "",
			t:      "2023-01-02T23:59:59.999999999+01:00",
			expect: 999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := NewFromTimeUsing(tTime, TotalNanoSeconds)
			if beats := newT.Beats(); beats != tt.expect {
				t.Errorf("expect %s to be @%d not @%d",
					tt.t,
					tt.expect,
					beats,
				)
			}
		})
	}
}
