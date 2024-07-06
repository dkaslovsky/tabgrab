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

// Environment variables
const (
	envVarPrefix      = "PREFIX"
	envVarBrowser     = "BROWSER"
	envVarBrowserArgs = "BROWSER_ARGS"
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
	fs.StringVar(&cFlags.browser, "browser", setStringFlagDefault(defaultBrowser, envVarBrowser), "browser name")
	fs.IntVar(&cFlags.maxTabs, "max", defaultMaxTabs, "maximum number of tabs")
	fs.StringVar(&cFlags.prefix, "prefix", setStringFlagDefault(defaultPrefix, envVarPrefix), "optional prefix for each URL")
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
	browserApp, validBrowser := browserApplications[strings.ToLower(cFlags.browser)]
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

	// Fall back on env var for prefix
	opts.prefix = cFlags.prefix
	if opts.prefix == "" {
		opts.prefix = os.Getenv(getEnvVarName(envVarPrefix))
	}

	opts.clipboard = cFlags.clipboard
	opts.verbose = cFlags.verbose

	return opts, nil
}

func setStringFlagDefault(defaultVal string, envVar string) string {
	if val := os.Getenv(getEnvVarName(envVar)); val != "" {
		return val
	}
	return defaultVal
}

func getEnvVarName(e string) string {
	return fmt.Sprintf("%s_%s", strings.ToUpper(appName), e)
}
