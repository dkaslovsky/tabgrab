package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

// Common flags
var (
	flagBrowser string
	flagMaxTabs int
	flagPrefix  string
)

func attachCommonFlags(fs *flag.FlagSet) {
	fs.StringVar(&flagBrowser, "browser", defaultBrowser, "browser name")
	fs.IntVar(&flagMaxTabs, "max", defaultMaxTabs, "maximum number of tabs")
	fs.StringVar(&flagPrefix, "prefix", defaultPrefix, "optional prefix for each URL")
}

type options struct {
	browserApp *browserApplication
	maxTabs    int
	prefix     string
}

func parseFlags(fs *flag.FlagSet, args []string, o *options) error {
	fs.Parse(args)

	// Set browser application
	browserApp, validBrowser := browserApplications[strings.ToLower(flagBrowser)]
	if !validBrowser {
		names := []string{}
		for name := range browserApplications {
			names = append(names, name)
		}
		sort.Strings(names)
		return fmt.Errorf("browser must be one of %v", names)
	}
	o.browserApp = browserApp

	// Set max tabs
	if flagMaxTabs <= 0 || flagMaxTabs > defaultMaxTabs {
		return fmt.Errorf("maximum tabs must be in the range [1, %d]", defaultMaxTabs)
	}
	o.maxTabs = flagMaxTabs

	// Set prefix
	o.prefix = flagPrefix

	return nil
}
