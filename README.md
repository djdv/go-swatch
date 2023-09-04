# go-swatch

Utility functions to convert Go `time`'s into [Swatch Internet Time](https://en.wikipedia.org/wiki/Swatch_Internet_Time).  
A decimal time relative to the Biel, Switzerland timezone, at a ratio of 1,000 ".beats" per 24 hour day.


As a command `cmd/swatch-time`:
```
$ swatch-time -h

Usage of swatch-time:
  -d    print date as well
  -p    use a more precise calculation method
  -r    use raw float format @000.000000
  -s    use Swatch standard format @000
(no flags defaults to centibeat format @000.00)
```

As a library:
```go
package main

import (
	"fmt"
	"time"

	"github.com/djdv/go-swatch"
)

func main() {
	// InternetTime values can be constructed
	// with the current time by calling [swatch.New].
	now := swatch.New()

	// The InternetTime type has the same method set
	// as a standard [time.Time].
	now.Day()
	now.Hour() // etc.

	// As well as some Swatch specific extensions.
	now.Beats()
	now.PreciseBeats()

	// It's possible to initialize an InternetTime
	// with some options, such as an existing time value,
	// and a more precise algorithm for resolving beats.
	const existingTime = "2006-02-15T02:57:08.000+01:00"
	standardTime, err := time.Parse(time.RFC3339, existingTime)
	if err != nil {
		panic(err)
	}
	swatchTime := swatch.New(
		swatch.WithTime(standardTime),
		swatch.WithAlgorithm(swatch.TotalNanoSeconds),
	)

	// Numerical beat values may be extracted from the time value.
	fmt.Println("Beats (int):", swatchTime.Beats())
	fmt.Println("High precision beats (float):", swatchTime.PreciseBeats())
	// Or formatted to strings using the same format tokens
	// from the [time] pkg, along with some swatch specific tokens.
	fmt.Println("Formatted centibeats:", swatchTime.Format(swatch.CentiBeats))
	fmt.Println("Formatted standard+swatch:", swatchTime.Format(time.DateOnly+swatch.Beats))
}
```

Output:
```
Beats (int): 123
High precision beats (float): 123.009259
Formatted centibeats: @123
Formatted standard+swatch: 2006-02-15@123
```
