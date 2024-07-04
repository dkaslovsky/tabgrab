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
	saveCmd    = flag.NewFlagSet("save", flag.ExitOnError)
	restoreCmd = flag.NewFlagSet("restore", flag.ExitOnError)
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

	case saveCmd.Name():
		attachCommonFlags(saveCmd)

		opts, err := parseSaveFlags(saveCmd, args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		err = saveTabs(opts)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case restoreCmd.Name():
		// fmt.Println("not implemented yet")

		attachCommonFlags(restoreCmd)

		opts, err := parseRestoreFlags(restoreCmd, args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		err = restoreTabs(opts)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

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
