package swatch_test

import (
	"math"
	"testing"
	"time"

	"github.com/djdv/go-swatch"
)

func TestInternetTime(t *testing.T) {
	t.Parallel()
	t.Run("constructor", constructor)
	t.Run("formatting", format)
	t.Run("behavior", behavior)
}

func constructor(t *testing.T) {
	t.Parallel()
	t.Run("new", newTime)
	t.Run("options", options)
}

func options(t *testing.T) {
	t.Parallel()
	t.Run("WithTime", fromTime)
}

func behavior(t *testing.T) {
	t.Parallel()
	t.Run("difference", durationDifference)
	t.Run("calculations", calculations)
}

func calculations(t *testing.T) {
	t.Parallel()
	t.Run("seconds", seconds)
	t.Run("nanoseconds", nanoseconds)
}

func seconds(t *testing.T) {
	t.Parallel()
	t.Run("standard", totalSecondsBeats)
	t.Run("precise", totalSecondsBeatsPrecise)
}

func nanoseconds(t *testing.T) {
	t.Parallel()
	t.Run("standard", totalNanosecondsBeats)
	t.Run("precise", totalNanosecondsBeatsPrecise)
}

func format(t *testing.T) {
	t.Parallel()
	t.Run("layout", timeLayout)
	t.Run("string", timeString)
	t.Run("combined", timeCombined)
}

func newTime(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected uninitialised InternetTime not to throw deref")
		}
	}()

	newT := swatch.New()
	if newT == nil {
		t.Errorf("expected New to return InternetTime")
		return
	}

	// Both of these functionally equivalent - just a sanity check
	_ = newT.UnixNano()
	_ = (&swatch.InternetTime{}).UnixNano()
}

func fromTime(t *testing.T) {
	t.Parallel()
	var (
		now  = time.Now()
		newT = swatch.New(swatch.WithTime(now))
	)
	if newT == nil {
		t.Errorf("expected NewFromTime to return InternetTime")
		return
	}

	if now.UnixNano() != newT.UnixNano() {
		t.Errorf("expected UnixNano of InternetTime and time.Time to match")
	}
}

func durationDifference(t *testing.T) {
	t.Parallel()
	// When given two dates exactly a beat apart, the beats are indeed 1 beat difference
	t1, err := time.Parse(time.RFC3339, "2006-02-15T12:00:00.000+01:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	t2, err := time.Parse(time.RFC3339, "2006-02-15T12:01:26.4000+01:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	var (
		i1 = swatch.New(swatch.WithTime(t1))
		i2 = swatch.New(swatch.WithTime(t2))

		a = roundDownFloat(i1.PreciseBeats(), 0)
		b = roundDownFloat(i2.PreciseBeats(), 0)
	)

	// Just to test sanity of GetTime
	if t1.Format("2006-01-02") != i1.Time.Format("2006-01-02") {
		t.Errorf("expected t1 date to be the same as swatchTime date")
	}

	if (b - a) < 1 {
		t.Errorf("expected b to be exactly 1 increment higher than a. Got a: %f b: %f", a, b)
	}
}

func totalSecondsBeats(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := swatch.New(swatch.WithTime(tTime))
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

func totalSecondsBeatsPrecise(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := swatch.New(swatch.WithTime(tTime))
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

func totalNanosecondsBeats(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := swatch.New(
				swatch.WithTime(tTime),
				swatch.WithAlgorithm(swatch.TotalNanoSeconds),
			)
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

func totalNanosecondsBeatsPrecise(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tTime, err := time.Parse(time.RFC3339, tt.t)
			if err != nil {
				t.Fatalf("error parsing test time: %s", err)
			}

			newT := swatch.New(
				swatch.WithTime(tTime),
				swatch.WithAlgorithm(swatch.TotalNanoSeconds),
			)
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

func timeLayout(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		format        string
		expectedValue string
		t             string
	}{
		{
			name:          "Swatch time format",
			format:        swatch.Beats,
			expectedValue: "@91",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Deci time format",
			format:        swatch.DeciBeats,
			expectedValue: "@91.2",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Centi time format",
			format:        swatch.CentiBeats,
			expectedValue: "@91.29",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Mili time format",
			format:        swatch.MilliBeats,
			expectedValue: "@91.296",
			t:             "2023-01-02T11:11:28+10:00",
		},
		{
			name:          "Micro time format",
			format:        swatch.MicroBeats,
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

func timeString(t *testing.T) {
	t.Parallel()
	tTime, err := time.Parse(time.RFC3339, "2023-01-02T11:11:28+10:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	newT := swatch.New(swatch.WithTime(tTime))

	if s := newT.String(); s != "@91" {
		t.Errorf("output of InternetTime String() unexpected: %s", s)
	}
}

func timeCombined(t *testing.T) {
	t.Parallel()
	tTime, err := time.Parse(time.RFC3339, "2023-01-02T11:11:28+10:00")
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}

	s := swatch.New(swatch.WithTime(tTime))
	if f := s.Format("2006-01-02 " + swatch.Beats); f != "2023-01-02 @91" {
		t.Errorf("Failed to mix formating, got %s", f)
	}
}

func equalWithTolerance(a, b float64) bool {
	tolerance := 0.000001 // Max length of Swatch Time
	diff := math.Abs(a - b)
	return diff < tolerance
}

func roundDownFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Floor(val*ratio) / ratio
}
