// Command swatch-time prints the current Swatch Internet Time in various .beat formats.
// Centibeats by default, Swatch standard with -s, and the raw underlying value with -r.
package main

import (
	"flag"
	"fmt"

	"github.com/djdv/go-swatch"
)

func main() {
	standard := flag.Bool("s", false, "use Swatch standard format @000")
	raw := flag.Bool("r", false, "use raw float format @000.000000")

	flag.Parse()

	if *raw && *standard {
		fmt.Println("Use none or 1 command flags, never both")
		return
	}

	if *raw {
		fmt.Println(swatch.Now(swatch.Raw))
		return
	}

	if *standard {
		fmt.Println(swatch.Now(swatch.Swatch))
		return
	}

	fmt.Println(swatch.Now(swatch.Centi))
}
