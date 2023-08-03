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
		raw, standard, precise, date bool
	)
	const (
		rawName      = "r"
		standardName = "s"
	)
	flagSet.Usage = usage
	flagSet.BoolVar(&raw, rawName, false, "use raw float format @000.000000")
	flagSet.BoolVar(&standard, standardName, false, "use Swatch standard format @000")
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
		fmt.Fprintf(flagSet.Output(),
			"Cannot combine -%s and -%s flags.",
			rawName, standardName,
		)
		return
	}
	var layout string
	switch {
	case raw:
		layout = swatch.MicroBeats
	case standard:
		layout = swatch.Beats
	default:
		layout = swatch.CentiBeats
	}

	var options []swatch.Option
	if precise {
		options = []swatch.Option{
			swatch.WithAlgorithm(swatch.TotalNanoSeconds),
		}
	}

	now := swatch.New(options...)
	if date {
		layout = "2006-01-02" + layout
	}

	fmt.Println(now.Format(layout))
}
