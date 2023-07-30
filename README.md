# go-swatch

Utility functions to convert Go `time`'s into [Swatch Internet Time](https://en.wikipedia.org/wiki/Swatch_Internet_Time).  
A decimal time relative to the Biel, Switzerland timezone, at a ratio of 1,000 ".beats" per 24 hour day.


As a command `cmd/swatch-time`
```
$ swatch-time -h

Usage of swatch-time:
  -d    print date as well
  -p    use a more precise calculation method
  -r    use raw float format @000.000000
  -s    use Swatch standard format @000
(no flags defaults to centibeat format @000.00)
```

## Initialising
Begin by importing the library:
```
import (
	"time"

	swatch "github.com/djdv/go-swatch/pkg/v1"
)
```

Grab the current time as a swatch internet time:
```
it := swatch.Now()
```

Or pass in your own time object to get the swatch time at a particular time:
```
t1, err := time.Parse(time.RFC3339, "2006-02-15T12:00:00.000+01:00")
if err != nil {
	panic("error parsing time")
}

it := swatch.NewFromTime(t1)
```

Perhaps you'd like to use better precision with the power of nanoseconds since Unix Epoch?
```
// The default is swatch.TotalSeconds
it := swatch.NewUsing(swatch.TotalNanoSeconds)
```

## Getting raw values
Raw values come in exactly two formats:
```
it.Beats() // int representing @000
```

Or using a very precise decimal:
```
it.PreciseBeats() // float64 representing @000.000000
```

Note: It's a time so it's rounded down based on the e+14 exponent (precision 7)

## Formatting an internet time
Lets start by printing Beats:
```
fmt.Println(it.Format(swatch.Swatch))
```

Perhaps you'd like the beats in the format @000.00?
```
fmt.Println(it.Format(swatch.Centi))
```

Because InternetTime is just an extension of time.Time you can use regular formatting:
```
fmt.Println(it.Format("2006-01-02"+swatch.Mili)) // Prints in the format YYYY-MM-DD@xxx.xxx
```

## Take a look under-the-hood
Need to change the algorithm after creating the swatch internet time?
```
it.SetAlgorithm(swatch.TotalNanoSeconds)
```

Need to get the underlying time.Time?
```
it.GetTime()
```

## More information?
Check out the implementation of `cmd/swatch-time/swatch-time.go` or puruse the tests for each file in the library under `pkg/v1`