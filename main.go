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

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	flag.Usage = func() { usage() }
	flag.Parse()

	subCmd, args := os.Args[1], os.Args[2:]

	switch subCmd {

	case grabCmd.Name():
		if err := runGrabCmd(grabCmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case tabCmd.Name(), tabCmdNameBackCompat:
		if err := runTabsCmd(tabCmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case versionCmd.Name():
		displayVersion()

	default:
		fmt.Printf("Error: unrecognized subcommand %s\n", subCmd)
		usage()
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

func usage() {
	fmt.Fprintf(os.Stderr, "%s: extract and restore URL tabs to and from the active browser window\n\n", appName)
	fmt.Fprint(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s:\t\t%s\n", grabCmdName, grabCmdDescription)
	fmt.Fprintf(os.Stderr, "  %s:\t\t%s\n", tabCmdName, tabCmdDescription)
	fmt.Fprintf(os.Stderr, "  %s:\t%s\n", versionCmdName, versionCmdDescription)
	fmt.Fprintf(os.Stderr, "\nRun `%s <subcommand> -help` for subcommand usage and flags", appName)
}
