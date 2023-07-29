package v1

import "strings"

func (t *internetTime) Format(layout string) string {
	// There's no "@" in time.Format src so it's safe to use as a delimeter
	// Replace in decending order of precision
	layout = strings.Replace(layout, Micro, t.format(Micro), 1)
	layout = strings.Replace(layout, Mili, t.format(Mili), 1)
	layout = strings.Replace(layout, Centi, t.format(Centi), 1)
	layout = strings.Replace(layout, Deci, t.format(Deci), 1)
	layout = strings.Replace(layout, Swatch, t.format(Swatch), 1)

	return t.Time.Format(layout)
}

func (t *internetTime) String() string {
	return t.Format(Swatch)
}

func (t *internetTime) format(layout string) string {
	return ""
}
