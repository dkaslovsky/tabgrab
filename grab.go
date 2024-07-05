package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type saveOptions struct {
	*commonOptions
}

func parseSaveFlags(fs *flag.FlagSet, args []string) (*saveOptions, error) {
	attachCommonFlags(fs)

	defaultUsage := fs.Usage
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "`%s` outputs the URL from each tab of the active browser window\n\n", grabCmdName)
		defaultUsage()
	}

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	commonOpts, err := parseCommonOptions()
	if err != nil {
		return nil, err
	}

	opts := &saveOptions{
		commonOptions: commonOpts,
	}
	return opts, nil
}

func saveTabs(opts *saveOptions) error {
	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Script to capture URL of tab i
	tabScript := "tell application \"" + opts.browserApp.cmdName + "\" to get {URL} of tab %d of window 1"

	for i := 0; i < opts.maxTabs; i++ {
		iTabScript := fmt.Sprintf(tabScript, i+1)

		// There is no intention that this implementation be secure so ignore the linter warning
		cmd := exec.Command("osascript", "-e", iTabScript) // #nosec
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if opts.verbose {
			log.Printf("executing: %s\n", cmd.String())
		}

		if err := cmd.Run(); err != nil {
			// Check stderr for clean exit on end-of-tabs error
			if stderr.Len() == 0 || isErrEndOfTabs(&stderr) {
				if opts.verbose {
					log.Printf("tab %d of window 1 not found, end of tabs", i+1)
				}
				break
			}
			// Return error if not end-of-tabs
			return fmt.Errorf("Error: %s\n%v\n", stderr.String(), err)
		}
	}

	tabs := strings.Split(stdout.String(), "\n")
	for _, tab := range tabs {
		if tab != "" {
			fmt.Printf("%s%s\n", opts.prefix, tab)
		}
	}

	return nil
}

// Error code indicating tab index out of range
var errCodeEndOfTabs = []byte("(-1719)\n")

func isErrEndOfTabs(buf *bytes.Buffer) bool {
	return bytes.HasSuffix(buf.Bytes(), errCodeEndOfTabs)
}
