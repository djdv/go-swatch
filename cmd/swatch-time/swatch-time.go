// Command swatch-time prints the current Swatch Internet Time in various .beat formats.
// Centibeats by default, Swatch standard with -s, and the raw underlying value with -r.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/djdv/go-swatch"
)

type settings struct {
	prefix, layout string
	options        []swatch.Option
}

func main() {
	var (
		settings = parseArgv()
		now      = swatch.New(settings.options...)
		stamp    = now.Format(settings.prefix + settings.layout)
	)
	fmt.Println(stamp)
}

func parseArgv() *settings {
	var (
		set = settings{
			layout: swatch.CentiBeats,
		}
		flagSet = newFlagSet(flag.ExitOnError)
		cmdName = flagSet.Name()
	)
	(&set).registerFlags(flagSet)
	flagSet.Usage = func() {
		output := flagSet.Output()
		fmt.Fprintf(output, "Usage of %s:\n", cmdName)
		flagSet.PrintDefaults()
		fmt.Fprint(output, "(no flags defaults to centibeat format @000.00)\n")
	}
	// We ignore this error because our
	// [flag.Flagset] is set to [flag.ExitOnError].
	_ = flagSet.Parse(os.Args[1:])
	if args := flagSet.Args(); len(args) > 0 {
		fmt.Fprintf(flagSet.Output(),
			"%s accepts no arguments but was passed: %s\n",
			cmdName, strings.Join(args, ", "),
		)
		flagSet.Usage()
		os.Exit(2)
	}
	return &set
}

func newFlagSet(eh flag.ErrorHandling) *flag.FlagSet {
	var (
		execName = filepath.Base(os.Args[0])
		cmdName  = strings.TrimSuffix(execName, filepath.Ext(execName))
	)
	return flag.NewFlagSet(cmdName, eh)
}

func (set *settings) registerFlags(flagSet *flag.FlagSet) {
	set.registerDateFlag(flagSet)
	set.registerFormatFlags(flagSet)
	set.registerPreciseFlag(flagSet)
}

func (set *settings) registerDateFlag(flagSet *flag.FlagSet) {
	const (
		name  = "d"
		usage = "print date as well"
	)
	flagSet.BoolFunc(name, usage, func(parameter string) error {
		useDate, err := strconv.ParseBool(parameter)
		if err != nil {
			return err
		}
		if useDate {
			set.prefix = time.DateOnly
		}
		return nil
	})
}

func (set *settings) registerFormatFlags(flagSet *flag.FlagSet) {
	const (
		rawName        = "r"
		rawUsage       = "use raw float format @000.000000"
		rawLayout      = swatch.MicroBeats
		standardName   = "s"
		standardUsage  = "use Swatch standard format @000"
		standardLayout = swatch.Beats
	)
	var (
		rawFlag, standardFlag bool
		parseFormatFlag       = func(parameter, layout string) error {
			if rawFlag && standardFlag {
				return fmt.Errorf(
					"cannot combine -%s and -%s flags",
					rawName, standardName,
				)
			}
			useLayout, err := strconv.ParseBool(parameter)
			if err != nil {
				return err
			}
			if useLayout {
				set.layout = layout
			}
			return nil
		}
	)
	flagSet.BoolFunc(rawName, rawUsage, func(parameter string) error {
		rawFlag = true
		return parseFormatFlag(parameter, rawLayout)
	})
	flagSet.BoolFunc(standardName, standardUsage, func(parameter string) error {
		standardFlag = true
		return parseFormatFlag(parameter, standardLayout)
	})
}

func (set *settings) registerPreciseFlag(flagSet *flag.FlagSet) {
	const (
		name  = "p"
		usage = "use a more precise calculation algorithm"
	)
	flagSet.BoolFunc(name, usage, func(parameter string) error {
		usePrecise, err := strconv.ParseBool(parameter)
		if err != nil {
			return err
		}
		if usePrecise {
			set.options = append(
				set.options,
				swatch.WithAlgorithm(swatch.TotalNanoSeconds),
			)
		}
		return nil
	})
}
