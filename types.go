package swatch

import "time"

type (
	Algorithm int
	Format    int
)

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

	swatchFormat = "@xxx"
	deciFormat   = "@xxx.x"
	centiFormat  = "@xxx.xx"
	miliFormat   = "@xxx.xxx"
	microFormat  = "@xxx.xxxxxx"
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
		return swatchFormat
	case Deci:
		return deciFormat
	case Centi:
		return centiFormat
	case Mili:
		return miliFormat
	case Micro:
		return microFormat
	default:
		return ""
	}
}
