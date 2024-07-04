package main

import (
	"flag"
	"fmt"
	"os"
)

// Application appName
const appName = "tabgrab"

// Version is set with ldflags
var version string

// Defaults
const (
	defaultBrowser = browserNameBrave
	defaultMaxTabs = 100
	defaultPrefix  = ""
)

// Subcommands
var (
	grabCmd    = flag.NewFlagSet("grab", flag.ExitOnError)
	openCmd    = flag.NewFlagSet("open", flag.ExitOnError)
	versionCmd = flag.NewFlagSet("version", flag.ExitOnError)
	helpCmd    = flag.NewFlagSet("help", flag.ExitOnError)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("PLACEHOLDER ERROR MESSAGE FOR MISSING SUBCMD")
		displayHelp()
		os.Exit(1)
	}

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {

	case grabCmd.Name():
		attachCommonFlags(grabCmd)

		opts, err := parseGrabFlags(grabCmd, args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		err = grabTabs(opts)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case openCmd.Name():
		fmt.Println("not implemented yet")

	case versionCmd.Name():
		displayVersion()

	case helpCmd.Name():
		displayHelp()

	default:
		fmt.Println("PLACEHOLDER ERROR MESSAGE FOR UNKNOWN SUBCMD")
		displayHelp()
		os.Exit(1)

	}
}

func displayVersion() {
	vStr := "%s: version %s\n"
	if version == "" {
		fmt.Printf(vStr, appName, "(development)")
		return
	}
	fmt.Printf(vStr, appName, version)
}

func displayHelp() {
	fmt.Println("PLACEHOLDER FOR HELP")
}
