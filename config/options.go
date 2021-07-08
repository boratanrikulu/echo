package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var ErrNoAddr = errors.New("server address is not defined")
var version = "v0.1.2"
var Banner = `
               _
              | |
  ___    ___  | |__     ___
 / _ \  / __| | '_ \   / _ \
|  __/ | (__  | | | | | (_) |
 \___|  \___| |_| |_|  \___/
                       ` + version

var usage = Banner + "\n" + `
Usage: echo [options]
Options:
	-a,  --address,   which interface and port will be used  *[example: ":1337"]
	-n,  --no-banner
	-v,  --verbose
	-h,  --help

* means "must be set".`

func printErrorAndDie(err error) {
	color.Red(err.Error())
	os.Exit(1)
}

func printHelpAndExit() {
	fmt.Println(usage)
	os.Exit(0)
}

type Options struct {
	Addr     string
	NoBanner bool
	Verbose  bool
	ShowHelp bool
}

// Configure sets options for the server.
func Configure(fs *flag.FlagSet, args []string) *Options {
	opts := &Options{}
	fs.Usage = func() {
		fmt.Println(usage)
	}

	fs.StringVar(&opts.Addr, "a", "", "")
	fs.StringVar(&opts.Addr, "address", "", "")
	fs.BoolVar(&opts.NoBanner, "n", false, "")
	fs.BoolVar(&opts.NoBanner, "no-banner", false, "")
	fs.BoolVar(&opts.Verbose, "v", false, "")
	fs.BoolVar(&opts.Verbose, "verbose", false, "")
	fs.BoolVar(&opts.ShowHelp, "h", false, "")
	fs.BoolVar(&opts.ShowHelp, "help", false, "")

	if err := fs.Parse(args); err != nil {
		printErrorAndDie(err)
	}
	if opts.ShowHelp {
		printHelpAndExit()
	}
	if opts.Addr == "" {
		printErrorAndDie(ErrNoAddr)
	}

	return opts
}
