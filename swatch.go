// Package swatch converts standard Go times to Swatch Internet Times.
// A decimal time relative to the Biel, Switzerland timezone, at a ratio of 1000 ".beats" per 24 hour day.
package swatch

import (
	"fmt"
	"time"
)

// Swatch Internet
type FormatStandard int

const (
	beatsPerSecond = 86.4

	// @000
	Swatch FormatStandard = iota
	// @000.00
	Centi
	// @000.000000
	Raw
)

// Beats formats a given time of day as a standard decimal amount of Swatch Beats.
func Beats(t time.Time) int {
	return int(RawBeats(t))
}

// CentiBeats formats a given time of day to a hundredth of a Swatch Beat.
func CentiBeats(t time.Time) float64 {
	return RawBeats(t) * 100 / 100
}

// RawBeats formats a given time of day to a b64 representation of a Swatch Beat.
func RawBeats(t time.Time) float64 {
	// Because swatch doesn't observe daylight savings, using CET will cause an error
	biel := time.FixedZone("UTC+1", 1*60*60) // 1 left in to demonstrate calculation

	tBiel := t.In(biel)
	// calculate seconds past midnight, divided it by beatsPerSecond
	return (float64(tBiel.Hour()*3600+tBiel.Minute()*60+tBiel.Second()) / beatsPerSecond)
}

// Now formats the current time in Internet Time.
func Now(standard FormatStandard) string {
	now := time.Now()

	switch standard {
	default:
		fallthrough
	case Swatch:
		return fmt.Sprintf("@%03d", Beats(now))
	case Centi:
		return fmt.Sprintf("@%06.2f", CentiBeats(now))
	case Raw:
		return fmt.Sprintf("@%06f", RawBeats(now))
	}
}
