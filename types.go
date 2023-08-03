package swatch

import "time"

type Algorithm int

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

const (
	TotalSeconds Algorithm = iota
	TotalNanoSeconds
)

type InternetTime struct {
	time.Time
	Algorithm
}
