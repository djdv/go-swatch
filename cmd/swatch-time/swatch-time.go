// Command swatch-time prints the current Swatch Internet Time in various .beat formats.
// Centibeats by default, Swatch standard with -s, and the raw underlying value with -r.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	swatch "github.com/djdv/go-swatch/pkg/v1"
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
		raw, standard, precise, date bool
	)
	flagSet.Usage = usage
	flagSet.BoolVar(&raw, "r", false, "use raw float format @000.000000")
	flagSet.BoolVar(&standard, "s", false, "use Swatch standard format @000")
	flagSet.BoolVar(&precise, "p", false, "use a more precise calculation method")
	flagSet.BoolVar(&date, "d", false, "print date as well")

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
		fmt.Fprint(flagSet.Output(), "Use none or 1 of the command flags, never combined.")
		return
	}
	var format swatch.Format
	switch {
	case raw:
		format = swatch.Micro
	case standard:
		format = swatch.Swatch
	default:
		format = swatch.Centi
	}

	s := swatch.Now()
	fs := format.String()
	if precise {
		s.SetAlgorithm(swatch.TotalNanoSeconds)
	}

	if date {
		fs = "2006-01-02" + fs
	}

	fmt.Println(s.Format(fs))
}
