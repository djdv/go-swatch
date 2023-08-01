package swatch

import "time"

type Algorithm int
type Format int

type InternetTime interface {
	// Raw time values
	Beats() int
	PreciseBeats() float64

	// Undercovers
	SetAlgorithm(Algorithm)
	GetTime() time.Time

	// Strings & Formatting
	String() string
	Format(string) string
}

const (
	Swatch Format = iota
	Deci
	Centi
	Mili
	Micro
)

const (
	TotalSeconds Algorithm = iota
	TotalNanoSeconds
)

type internetTime struct {
	time.Time
	algorithm Algorithm
}

func (f Format) String() string {
	switch f {
	case Swatch:
		return "@xxx"
	case Deci:
		return "@xxx.x"
	case Centi:
		return "@xxx.xx"
	case Mili:
		return "@xxx.xxx"
	case Micro:
		return "@xxx.xxxxxx"
	}

	return ""
}
