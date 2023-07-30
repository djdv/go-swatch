package v1

import "time"

type Algorithm int
type Format int

type InternetTime interface {
	// Raw time values
	Beats() int
	PreciseBeats() float64

	// Configuration
	SetAlgorithm(Algorithm)

	// Strings & Formatting
	String() string
	Format(string) string
}

const (
	Swatch Format = iota // 0
	Deci                 // 1
	Centi                // 2
	Mili                 // 3
	Micro                // 4

	TotalSeconds     Algorithm = iota // 0
	TotalNanoSeconds                  // 1
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
