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

// Subcommand names
const (
	saveCmdName    = "save"
	restoreCmdName = "restore"
	versionCmdName = "version"
	helpCmdName    = "help"
)

// Subcommands
var (
	saveCmd    = flag.NewFlagSet(saveCmdName, flag.ExitOnError)
	restoreCmd = flag.NewFlagSet(restoreCmdName, flag.ExitOnError)
	versionCmd = flag.NewFlagSet(versionCmdName, flag.ExitOnError)
	helpCmd    = flag.NewFlagSet(helpCmdName, flag.ExitOnError)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("PLACEHOLDER ERROR MESSAGE FOR MISSING SUBCMD")
		displayHelp()
		os.Exit(1)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", appName)
		fmt.Fprintf(os.Stderr, "PLACEHOLDER FOR COMMANDS: %v", []string{saveCmdName, restoreCmdName, versionCmdName, helpCmdName})
	}
	flag.Parse()

	cmd, args := os.Args[1], os.Args[2:]

	switch cmd {

	case saveCmd.Name():
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
