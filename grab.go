package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

type grabOptions struct {
	*commonOptions
}

func parseGrabFlags(fs *flag.FlagSet, args []string) (*grabOptions, error) {
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	commonOpts, err := setCommonOptions()
	if err != nil {
		return nil, err
	}

	opts := &grabOptions{
		commonOptions: commonOpts,
	}
	return opts, nil
}

func grabTabs(opts *grabOptions) error {
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
