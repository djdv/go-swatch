// Command swatch-time prints the current Swatch Internet Time in various .beat formats.
// Centibeats by default, Swatch standard with -s, and the raw underlying value with -r.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/djdv/go-swatch"
)

func main() {
	var (
		execName = filepath.Base(os.Args[0])
		cmdName  = strings.TrimSuffix(execName, filepath.Ext(execName))
		flagSet  = flag.NewFlagSet(cmdName, flag.ExitOnError)
		usage    = func() {
			output := flagSet.Output()
			fmt.Fprintf(output, "Usage of %s:\n", cmdName)
			flagSet.PrintDefaults()
			fmt.Fprint(output, "(no flags defaults to centibeat format @000.00)\n")
		}
		raw, standard bool
	)
	flagSet.Usage = usage
	flagSet.BoolVar(&raw, "r", false, "use raw float format @000.000000")
	flagSet.BoolVar(&standard, "s", false, "use Swatch standard format @000")

	if flagSet.Parse(os.Args[1:]) != nil {
		return
	}
	if args := flagSet.Args(); len(args) > 0 {
		fmt.Fprintf(flagSet.Output(),
			"%s accepts no arguments but was passed: %s\n",
			cmdName, strings.Join(args, ", "),
		)
		flagSet.Usage()
		return
	}

	if raw && standard {
		fmt.Println("Use none or 1 command flags, never both")
		return
	}

	if raw {
		fmt.Println(swatch.Now(swatch.Raw))
		return
	}

	if standard {
		fmt.Println(swatch.Now(swatch.Swatch))
		return
	}

	fmt.Println(swatch.Now(swatch.Centi))
}
