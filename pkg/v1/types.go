package v1

import "time"

type Algorithm = int

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
	Swatch = "@000"
	Deci   = "@000.0"
	Centi  = "@000.00"
	Mili   = "@000.000"
	Micro  = "@000.000000"

	TotalSeconds     Algorithm = iota // 0
	TotalNanoSeconds                  // 1
)

type internetTime struct {
	time.Time
	algorithm Algorithm
}
