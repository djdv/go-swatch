package swatch

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *InternetTime) Format(layout string) string {
	// There's no "@" in time.Format src so it's safe to use as a delimiter
	// Replace in descending order of precision
	return strings.NewReplacer(
		MicroBeats, t.format(MicroBeats),
		MilliBeats, t.format(MilliBeats),
		CentiBeats, t.format(CentiBeats),
		DeciBeats, t.format(DeciBeats),
		Beats, t.format(Beats),
	).Replace(t.Time.Format(layout))
}

func (t *InternetTime) String() string {
	return t.Format(Beats)
}

// Expects layout to only be one of the predefined format tokens.
func (t *InternetTime) format(layout string) string {
	switch layout {
	case Beats:
		return fmt.Sprintf("@%d", t.Beats())
	case DeciBeats, CentiBeats, MilliBeats:
		beats := roundDownFloat(t.PreciseBeats(), precisionOf(layout))
		return fmt.Sprintf("@%s", strconv.FormatFloat(beats, 'f', -1, 64))
	case MicroBeats:
		return fmt.Sprintf("@%s", strconv.FormatFloat(t.PreciseBeats(), 'f', -1, 64))
	default:
		return ""
	}
}
