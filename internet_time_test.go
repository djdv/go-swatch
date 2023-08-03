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
	t.Run("new", swatchNew)
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
	t.Run("combined", dateAndTime)
}

func swatchNew(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("uninitialised InternetTime panicked: %v", r)
		}
	}()

	newT := newSwatchTime(t)

	_ = newT.UnixNano()
	_ = (&swatch.InternetTime{}).UnixNano()
}

func fromTime(t *testing.T) {
	t.Parallel()
	var (
		now  = time.Now()
		newT = newSwatchTime(t, swatch.WithTime(now))
		got  = newT.UnixNano()
		want = now.UnixNano()
	)
	if got != want {
		t.Errorf("UnixNano time mismatch"+
			"\n\tgot: %d"+
			"\n\twant: %d",
			got, want,
		)
	}
}

func durationDifference(t *testing.T) {
	t.Parallel()
	// Validate that two dates exactly a beat apart, are indeed 1 beat difference.
	var (
		t1 = parseStandardTime(t, time.RFC3339, "2006-02-15T12:00:00.000+01:00")
		t2 = parseStandardTime(t, time.RFC3339, "2006-02-15T12:01:26.4000+01:00")

		i1 = newSwatchTime(t, swatch.WithTime(t1))
		i2 = newSwatchTime(t, swatch.WithTime(t2))

		beats1 = roundDownFloat(i1.PreciseBeats(), 0)
		beats2 = roundDownFloat(i2.PreciseBeats(), 0)

		got  = t1.Format("2006-01-02")
		want = i1.Time.Format("2006-01-02")
	)

	if got != want {
		t.Errorf("standard time and internet time mismatch"+
			"\n\tgot: %s"+
			"\n\twant: %s",
			got, want,
		)
	}

	if (beats2 - beats1) != 1 {
		t.Errorf("duration difference should be exactly 1 beat apart"+
			"\n\tgot: %f"+
			"\n\twant: %f",
			beats1, beats2,
		)
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
			var (
				tTime = parseStandardTime(t, time.RFC3339, tt.t)
				newT  = newSwatchTime(t, swatch.WithTime(tTime))
				got   = newT.Beats()
				want  = tt.expect
			)
			compareBeats(t, got, want, tt.t)
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
			var (
				tTime = parseStandardTime(t, time.RFC3339, tt.t)
				newT  = newSwatchTime(t, swatch.WithTime(tTime))
				got   = newT.PreciseBeats()
				want  = tt.expect
			)
			compareBeats(t, got, want, tt.t)
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
			var (
				tTime = parseStandardTime(t, time.RFC3339, tt.t)
				newT  = swatch.New(
					swatch.WithTime(tTime),
					swatch.WithAlgorithm(swatch.TotalNanoSeconds),
				)
				got  = newT.Beats()
				want = tt.expect
			)
			compareBeats(t, got, want, tt.t)
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
			var (
				tTime = parseStandardTime(t, time.RFC3339, tt.t)
				newT  = swatch.New(
					swatch.WithTime(tTime),
					swatch.WithAlgorithm(swatch.TotalNanoSeconds),
				)
				got  = newT.PreciseBeats()
				want = tt.expect
			)
			compareBeats(t, got, want, tt.t)
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
			name:          "Milli time format",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				tTime = parseStandardTime(t, time.RFC3339, tt.t)
				newT  = newSwatchTime(t, swatch.WithTime(tTime))
				got   = newT.Format(tt.format)
				want  = tt.expectedValue
			)
			compareLayout(t, got, want, tTime)
		})
	}
}

func timeString(t *testing.T) {
	t.Parallel()
	var (
		tTime = parseStandardTime(t, time.RFC3339, "2023-01-02T11:11:28+10:00")
		newT  = newSwatchTime(t, swatch.WithTime(tTime))
		got   = newT.String()
		want  = "@91"
	)
	compareLayout(t, got, want, tTime)
}

func dateAndTime(t *testing.T) {
	t.Parallel()
	var (
		tTime = parseStandardTime(t, time.RFC3339, "2023-01-02T11:11:28+10:00")
		s     = newSwatchTime(t, swatch.WithTime(tTime))
		got   = s.Format("2006-01-02 " + swatch.Beats)
		want  = "2023-01-02 @91"
	)
	compareLayout(t, got, want, tTime)
}

func newSwatchTime(t *testing.T, options ...swatch.Option) *swatch.InternetTime {
	t.Helper()
	newT := swatch.New(options...)
	if newT == nil {
		t.Fatal("New did not returned a valid value")
	}
	return newT
}

func parseStandardTime(t *testing.T, layout, value string) time.Time {
	t.Helper()
	stdTime, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("error parsing test time: %s", err)
	}
	return stdTime
}

func compareBeats[T int | float64](t *testing.T, got, want T, formatted string) {
	t.Helper()
	var match bool
	// Workaround for: golang/go #45380.
	if _, precise := any((*T)(nil)).(*float64); precise {
		match = equalWithTolerance(float64(got), float64(want))
	} else {
		match = got == want
	}
	if !match {
		t.Errorf("beats mismatch for %s"+
			"\n\tgot: @%v"+
			"\n\twant: @%v",
			formatted,
			got, want,
		)
	}
}

func compareLayout(t *testing.T, got, want string, stdTime time.Time) {
	t.Helper()
	if got != want {
		t.Errorf("layout format mismatch for time %s"+
			"\n\tgot: %s"+
			"\n\twant: %s",
			stdTime,
			got, want,
		)
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
