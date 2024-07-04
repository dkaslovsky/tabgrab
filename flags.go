package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

type commonFlags struct {
	browser string
	maxTabs int
	prefix  string
}

var commonFlagValues commonFlags

func attachCommonFlags(fs *flag.FlagSet) {
	fs.StringVar(&commonFlagValues.browser, "browser", defaultBrowser, "browser name")
	fs.IntVar(&commonFlagValues.maxTabs, "max", defaultMaxTabs, "maximum number of tabs")
	fs.StringVar(&commonFlagValues.prefix, "prefix", defaultPrefix, "optional prefix for each URL")
}

type commonOptions struct {
	browserApp *browserApplication
	maxTabs    int
	prefix     string
}

func setCommonOptions() (*commonOptions, error) {
	opts := &commonOptions{}

	// Set browser application
	browserApp, validBrowser := browserApplications[strings.ToLower(commonFlagValues.browser)]
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
	if commonFlagValues.maxTabs <= 0 || commonFlagValues.maxTabs > defaultMaxTabs {
		return nil, fmt.Errorf("maximum tabs must be in the range [1, %d]", defaultMaxTabs)
	}
	opts.maxTabs = commonFlagValues.maxTabs

	// Set prefix
	opts.prefix = commonFlagValues.prefix

	return opts, nil
}
