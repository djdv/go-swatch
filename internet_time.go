package swatch

import (
	"math"
	"time"
)

const (
	defaultAlgorithm         = TotalSeconds
	secondsPerBeat           = 86.4
	nanoPerDay         int64 = 8.64e+13
	nanoPerBeat        int64 = nanoPerDay / 1000
	nanoPerHour        int64 = 3.6e+12
	maxSwatchPrecision int   = 6
)

type Option func(*InternetTime)

func New(options ...Option) *InternetTime {
	swatchTime := InternetTime{
		Algorithm: defaultAlgorithm,
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

func (t *InternetTime) Beats() int {
	return int(roundDownFloat(t.PreciseBeats(), 0))
}

func (t *InternetTime) PreciseBeats() float64 {
	return t.calculateBeats()
}

func (t *InternetTime) calculateBeats() float64 {
	n := float64(0)
	switch t.Algorithm {
	case TotalSeconds:
		n = totalSecondsAlgorithm(t.Time)
	case TotalNanoSeconds:
		n = totalNanosecondsAlgorithm(t.Time)
	}

	n = roundDownFloat(n, maxSwatchPrecision)

	return n
}

func totalSecondsAlgorithm(t time.Time) float64 {
	hourSeconds := t.Hour() * 3600
	minuteSeconds := t.Minute() * 60
	totalSeconds := float64(hourSeconds+minuteSeconds+t.Second()) / secondsPerBeat
	return totalSeconds
}

func totalNanosecondsAlgorithm(t time.Time) float64 {
	// Convert to UTC as that's what unix timestamp will be
	t = t.In(time.UTC)
	// t.UnixNano() + nanoPerHour accounts for internetTime being UTC+1
	nanoSecondsSinceYesterday := (t.UnixNano() + nanoPerHour) % nanoPerDay

	if nanoSecondsSinceYesterday == 0 {
		nanoSecondsSinceYesterday = nanoPerDay
	}

	totalSeconds := float64(nanoSecondsSinceYesterday) / float64(nanoPerBeat)
	return totalSeconds
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
