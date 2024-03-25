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

const (
	// Defaults
	defaultBrowser = browserNameBrave
	defaultMaxTabs = 100

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

	fmt.Printf("%s", stdout.String())
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
}

func parseFlags(o *options) error {
	// Flags with defaults
	browser := flag.String("browser", defaultBrowser, "browser name")
	maxTabs := flag.Int("max", defaultMaxTabs, "maximum number of tabs")
	flag.Parse()

	// Set browser application
	if browserApp, ok := browserApplications[strings.ToLower(*browser)]; ok {
		o.browserApp = browserApp
	} else {
		names := []string{}
		for name := range browserApplications {
			names = append(names, name)
		}
		sort.Strings(names)
		return fmt.Errorf("browser must be one of %v", names)
	}

	// Set max tabs
	if *maxTabs <= 0 || *maxTabs > defaultMaxTabs {
		return fmt.Errorf("maximum tabs must be in the range [1, %d]", defaultMaxTabs)
	}
	o.maxTabs = *maxTabs

	return nil
}

// Error code indicating tab index out of range
var errCodeEndOfTabs = []byte("(-1719)\n")

func isErrEndOfTabs(buf *bytes.Buffer) bool {
	return bytes.HasSuffix(buf.Bytes(), errCodeEndOfTabs)
}
