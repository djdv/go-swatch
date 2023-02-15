# go-swatch

Utility functions to convert Go `time`'s into [Swatch Internet Time](https://en.wikipedia.org/wiki/Swatch_Internet_Time).  
A decimal time relative to the Biel, Switzerland timezone, at a ratio of 1,000 ".beats" per 24 hour day.


As a command `cmd/swatch-time`
```
swatch-time -h
Usage of swatch-time:
  -r    use raw float format @000.000000
  -s    use Swatch standard format @000
(no flags defaults to centibeat format @000.00)
```

As a library
```go
// Get numerical representations:
var (
	now        = time.Now()
	beatsInt   = swatch.Beats(now)
	centiFloat = swatch.CentiBeats(now)
	rawFloat   = swatch.RawBeats(now)
)
// Get current time as beats formatted string:
// Swatch: @000
// Centi:  @000.00
// Raw:    @000.000000
swatch.Now(swatch.Swatch)
```
