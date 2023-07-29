package v1

import (
	"math"
	"time"
)

var (
	defaultAlgorithm       = TotalSeconds
	secondsPerBeat         = 86.4
	nanoPerDay       int64 = 86.4e+13
	nanoPerBeat      int64 = nanoPerDay / 1000
)

func New() *internetTime {
	return &internetTime{getUtcTime(time.Now()), defaultAlgorithm}
}

func NewFromTime(t time.Time) *internetTime {
	return &internetTime{getUtcTime(t), defaultAlgorithm}
}

func NewUsing(algo Algorithm) *internetTime {
	i := New()
	i.SetAlgorithm(algo)
	return i
}

func NewFromTimeUsing(t time.Time, algo Algorithm) *internetTime {
	i := NewFromTime(t)
	i.SetAlgorithm(algo)
	return i
}

func (t *internetTime) SetAlgorithm(algo Algorithm) {
	t.algorithm = algo
}

func (t *internetTime) Beats() int {
	return int(roundDownFloat(t.PreciseBeats(), 0))
}

func (t *internetTime) PreciseBeats() float64 {
	return t.calculateBeats()
}

func (t *internetTime) calculateBeats() float64 {
	switch t.algorithm {
	case TotalSeconds:
		return totalSecondsAlgorithm(t.Time)
	case TotalNanoSeconds:
		return totalNanosecondsAlgorithm(t.Time)
	}

	return 0
}

func totalSecondsAlgorithm(t time.Time) float64 {
	return float64(t.Hour()*3600+t.Minute()*60+t.Second()) / secondsPerBeat
}

func totalNanosecondsAlgorithm(t time.Time) float64 {
	nanoSecondsSinceYesterday := t.UnixNano() % nanoPerDay
	if nanoSecondsSinceYesterday == 0 {
		return 1
	}

	return float64(nanoSecondsSinceYesterday) / float64(nanoPerBeat)
}

func roundDownFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Floor(val*ratio) / ratio
}

func getUtcTime(t time.Time) time.Time {
	// Because swatch doesn't observe daylight savings, using CET will cause an error
	biel := time.FixedZone("UTC+1", 1*60*60) // 1 left in to demonstrate calculation
	return t.In(biel)
}
