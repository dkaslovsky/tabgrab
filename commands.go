package main

import "flag"

// Subcommand names
const (
	grabCmdName          = "grab"
	tabCmdName           = "tabs"
	tabCmdNameBackCompat = "tab" // Backwards compatibility with old command name
	closeCmdName         = "close"
	versionCmdName       = "version"
)

// Subcommands
var (
	grabCmd    = flag.NewFlagSet(grabCmdName, flag.ExitOnError)
	tabCmd     = flag.NewFlagSet(tabCmdName, flag.ExitOnError)
	closeCmd   = flag.NewFlagSet(closeCmdName, flag.ExitOnError)
	versionCmd = flag.NewFlagSet(versionCmdName, flag.ExitOnError)
)

// Subcommand descriptions
const (
	grabCmdDescription    = "extracts the URL from each tab of the active browser window"
	tabCmdDescription     = "opens the provided URLs as tabs in a new browser window"
	closeCmdDescription   = "DESCRIPTION GOES HERE" // TODO: fix this
	versionCmdDescription = "displays application version information"
)
