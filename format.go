package swatch

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *internetTime) Format(layout string) string {
	// There's no "@" in time.Format src so it's safe to use as a delimiter
	// Replace in descending order of precision
	return strings.NewReplacer(
		microFormat, t.format(Micro),
		miliFormat, t.format(Mili),
		centiFormat, t.format(Centi),
		deciFormat, t.format(Deci),
		swatchFormat, t.format(Swatch),
	).Replace(t.Time.Format(layout))
}

func (t *internetTime) String() string {
	return t.Format(swatchFormat)
}

// Expects layout to only be
func (t *internetTime) format(layout Format) string {
	r := ""
	switch layout {
	case Swatch:
		r = fmt.Sprintf("@%d", t.Beats())
	case Deci, Centi, Mili:
		beats := roundDownFloat(t.PreciseBeats(), int(layout))
		r = fmt.Sprintf("@%s", strconv.FormatFloat(beats, 'f', -1, 64))
	case Micro:
		r = fmt.Sprintf("@%s", strconv.FormatFloat(t.PreciseBeats(), 'f', -1, 64))
	}

	return r
}
