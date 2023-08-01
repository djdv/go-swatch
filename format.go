package swatch

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *internetTime) Format(layout string) string {
	// There's no "@" in time.Format src so it's safe to use as a delimeter
	// Replace in decending order of precision
	layout = t.Time.Format(layout)

	layout = strings.Replace(layout, Micro.String(), t.format(Micro), 1)
	layout = strings.Replace(layout, Mili.String(), t.format(Mili), 1)
	layout = strings.Replace(layout, Centi.String(), t.format(Centi), 1)
	layout = strings.Replace(layout, Deci.String(), t.format(Deci), 1)
	layout = strings.Replace(layout, Swatch.String(), t.format(Swatch), 1)

	return layout
}

func (t *internetTime) String() string {
	return t.Format(Swatch.String())
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
