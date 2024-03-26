package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Application name
const name = "tabgrab"

// Version is set with ldflags
var version string

const (
	// Defaults
	defaultBrowser = browserNameBrave
	defaultMaxTabs = 100
	defaultPrefix  = ""

	// Browser names
	browserNameBrave  = "brave"
	browserNameChrome = "chrome"
	browserNameSafari = "safari"
)

// Browser applications by name
var browserApplications = map[string]*browserApplication{
	browserNameBrave:  {"Brave Browser"},
	browserNameChrome: {"Google Chrome"},
	browserNameSafari: {"Safari"},
}

func main() {
	opts := &options{}
	if err := parseFlags(opts); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if opts.version {
		displayVersion()
		os.Exit(0)
	}

	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Script to capture URL of tab i
	tabScript := "tell application \"" + opts.browserApp.String() + "\" to get {URL} of tab %d of window 1"

	for i := 0; i < opts.maxTabs; i++ {
		iTabScript := fmt.Sprintf(tabScript, i+1)

		// There is no intention that this implementation be secure so ignore the linter warning
		cmd := exec.Command("osascript", "-e", iTabScript) // #nosec
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			// Check stderr for clean exit on end-of-tabs error
			if stderr.Len() == 0 || isErrEndOfTabs(&stderr) {
				break
			}
			// Report error if not end-of-tabs
			fmt.Printf("Error: %s\n%v\n", stderr.String(), err)
			os.Exit(1)
		}
	}

	tabs := strings.Split(stdout.String(), "\n")
	for _, tab := range tabs {
		if tab != "" {
			fmt.Printf("%s%s\n", opts.prefix, tab)
		}
	}
}

type browserApplication struct {
	appName string
}

func (b browserApplication) String() string {
	return b.appName
}

type options struct {
	browserApp *browserApplication
	maxTabs    int
	prefix     string
	version    bool
}

func parseFlags(o *options) error {
	// Flags with defaults
	var (
		browser = flag.String("browser", defaultBrowser, "browser name")
		maxTabs = flag.Int("max", defaultMaxTabs, "maximum number of tabs")
		prefix  = flag.String("prefix", defaultPrefix, "optional prefix to attach to each tab's URL")
		version = flag.Bool("version", false, "display version")
	)

	flag.Parse()

	if *version {
		o.version = true
		return nil
	}

	// Set browser application
	browserApp, validBrowser := browserApplications[strings.ToLower(*browser)]
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
	if *maxTabs <= 0 || *maxTabs > defaultMaxTabs {
		return fmt.Errorf("maximum tabs must be in the range [1, %d]", defaultMaxTabs)
	}
	o.maxTabs = *maxTabs

	// Set prefix
	o.prefix = *prefix

	return nil
}

// Error code indicating tab index out of range
var errCodeEndOfTabs = []byte("(-1719)\n")

func isErrEndOfTabs(buf *bytes.Buffer) bool {
	return bytes.HasSuffix(buf.Bytes(), errCodeEndOfTabs)
}

func displayVersion() {
	vStr := "%s: version %s\n"
	if version == "" {
		fmt.Printf(vStr, name, "(development)")
		return
	}
	fmt.Printf(vStr, name, "v"+version)
}
