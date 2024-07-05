package main

import "flag"

// Subcommand names
const (
	grabCmdName    = "grab"
	tabCmdName     = "tab"
	versionCmdName = "version"
)

// Subcommands
var (
	grabCmd    = flag.NewFlagSet(grabCmdName, flag.ExitOnError)
	tabCmd     = flag.NewFlagSet(tabCmdName, flag.ExitOnError)
	versionCmd = flag.NewFlagSet(versionCmdName, flag.ExitOnError)
)

// Subcommand descriptions
const (
	grabCmdDescription    = "extracts the URL from each tab of the active browser window"
	tabCmdDescription     = "opens the provided URLs as tabs in a new browser window"
	versionCmdDescription = "displays application version information"
)
