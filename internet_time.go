package swatch

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type (
	// Algorithm defines the method used
	// when calculating the beats of a
	// time value.
	Algorithm    func(time.Time) float64
	InternetTime struct {
		time.Time
		Algorithm
	}
	Option func(*InternetTime)
)

const (
	// Beats is a format token which gets replaced
	// with the standard Swatch Internet Time .beats layout.
	Beats = "@xxx"
	// DeciBeats is a format token which gets replaced
	// with a higher precision .beats layout.
	DeciBeats = "@xxx.x"
	// CentiBeats is a format token which gets replaced
	// with a higher precision .beats layout.
	// Also sometimes called "sub-beats".
	CentiBeats = "@xxx.xx"
	// MilliBeats is a format token which gets replaced
	// with a higher precision .beats layout.
	MilliBeats = "@xxx.xxx"
	// MicroBeats is a format token which gets replaced
	// with a higher precision .beats layout.
	MicroBeats = "@xxx.xxxxxx"
)

// TotalSeconds is an [Algorithm] which uses
// a time resolution of seconds to calculate
// beats.
func TotalSeconds(t time.Time) float64 {
	const secondsPerBeat = 86.4
	var (
		hourSeconds   = t.Hour() * 3600
		minuteSeconds = t.Minute() * 60
		seconds       = t.Second()
		sum           = float64(hourSeconds + minuteSeconds + seconds)
	)
	return sum / secondsPerBeat
}

// TotalNanoSeconds is an [Algorithm] which uses
// a time resolution of nanoseconds to calculate
// beats.
func TotalNanoSeconds(t time.Time) float64 {
	const (
		nanoPerDay  int64 = 8.64e+13
		nanoPerBeat int64 = nanoPerDay / 1000
		nanoPerHour int64 = 3.6e+12
	)
	var (
		utc  = t.In(time.UTC) // Normalize to UTC.
		nano = utc.UnixNano()
		// Account for Internet Time being UTC+1.
		sinceYesterday = (nano + nanoPerHour) % nanoPerDay
	)
	if sinceYesterday == 0 {
		sinceYesterday = nanoPerDay
	}
	return float64(sinceYesterday) / float64(nanoPerBeat)
}

func New(options ...Option) *InternetTime {
	swatchTime := InternetTime{
		Algorithm: TotalSeconds,
	}
	for _, apply := range options {
		apply(&swatchTime)
	}
	if swatchTime.Time.IsZero() {
		swatchTime.Time = getUtcTime(time.Now())
	}
	return &swatchTime
}

func WithAlgorithm(algorithm Algorithm) Option {
	return func(it *InternetTime) {
		it.Algorithm = algorithm
	}
}

func WithTime(t time.Time) Option {
	return func(it *InternetTime) {
		it.Time = getUtcTime(t)
	}
}

func (t *InternetTime) Format(layout string) string {
	// There's no "@" in time.Format src so it's safe to use as a delimiter
	// Replace in descending order of precision
	return strings.NewReplacer(
		MicroBeats, t.format(MicroBeats),
		MilliBeats, t.format(MilliBeats),
		CentiBeats, t.format(CentiBeats),
		DeciBeats, t.format(DeciBeats),
		Beats, t.format(Beats),
	).Replace(t.Time.Format(layout))
}

// Expects layout to only be one of the predefined format tokens.
func (t *InternetTime) format(layout string) string {
	switch layout {
	case Beats:
		return fmt.Sprintf("@%d", t.Beats())
	case DeciBeats, CentiBeats, MilliBeats:
		beats := roundDownFloat(t.PreciseBeats(), precisionOf(layout))
		return fmt.Sprintf("@%s", strconv.FormatFloat(beats, 'f', -1, 64))
	case MicroBeats:
		return fmt.Sprintf("@%s", strconv.FormatFloat(t.PreciseBeats(), 'f', -1, 64))
	default:
		return ""
	}
}

func (t *InternetTime) String() string {
	return t.Format(Beats)
}

func (t *InternetTime) Beats() int {
	return int(roundDownFloat(t.PreciseBeats(), 0))
}

func (t *InternetTime) PreciseBeats() float64 {
	return t.calculateBeats()
}

func (t *InternetTime) calculateBeats() float64 {
	const maxSwatchPrecision int = 6
	beats := t.Algorithm(t.Time)
	return roundDownFloat(beats, maxSwatchPrecision)
}

func precisionOf(format string) int {
	switch format {
	case DeciBeats:
		return 1
	case CentiBeats:
		return 2
	case MilliBeats:
		return 3
	default:
		return 0
	}
}

func roundDownFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Floor(val*ratio) / ratio
}

func getUtcTime(t time.Time) time.Time {
	t = t.UTC()
	// Because swatch doesn't observe daylight savings, using CET will cause an error
	biel := time.FixedZone("UTC+1", 1*60*60) // 1 left in to demonstrate calculation
	return t.In(biel)
}
