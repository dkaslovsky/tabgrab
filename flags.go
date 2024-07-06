package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Defaults
const (
	defaultBrowser = browserNameChrome
	defaultMaxTabs = 100
	defaultPrefix  = ""
)

type commonFlags struct {
	browser   string
	maxTabs   int
	prefix    string
	clipboard bool
	verbose   bool
}

// Global instance of the common flag values struct
var cFlags commonFlags

func attachCommonFlags(fs *flag.FlagSet) {
	fs.StringVar(&cFlags.browser, "browser", defaultBrowser, "browser name")
	fs.IntVar(&cFlags.maxTabs, "max", defaultMaxTabs, "maximum number of tabs")
	fs.StringVar(&cFlags.prefix, "prefix", defaultPrefix, "optional prefix for each URL")
	fs.BoolVar(&cFlags.clipboard, "clipboard", false, "use clipboard for input/output")
	fs.BoolVar(&cFlags.verbose, "verbose", false, "enable verbose output")
}

type commonOptions struct {
	browserApp *browserApplication
	maxTabs    int
	prefix     string
	clipboard  bool
	verbose    bool
}

func parseCommonOptions() (*commonOptions, error) {
	opts := &commonOptions{}

	// Set browser application
	browser := cFlags.browser
	if browserOverride := os.Getenv("TABGRAB_BROWSER"); browserOverride != "" {
		browser = browserOverride
	}
	browserApp, validBrowser := browserApplications[strings.ToLower(browser)]
	if !validBrowser {
		names := []string{}
		for name := range browserApplications {
			names = append(names, name)
		}
		sort.Strings(names)
		return nil, fmt.Errorf("browser must be one of %v", names)
	}
	opts.browserApp = browserApp

	// Set max tabs
	if cFlags.maxTabs <= 0 || cFlags.maxTabs > defaultMaxTabs {
		return nil, fmt.Errorf("maximum tabs must be in the range [1, %d]", defaultMaxTabs)
	}
	opts.maxTabs = cFlags.maxTabs

	opts.prefix = cFlags.prefix
	opts.clipboard = cFlags.clipboard
	opts.verbose = cFlags.verbose

	return opts, nil
}
