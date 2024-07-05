package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func runGrabCmd(cmd *flag.FlagSet, args []string) error {
	opts, err := parseGrabFlags(cmd, args)
	if err != nil {
		return err
	}

	err = grabTabs(opts)
	if err != nil {
		return err
	}

	return nil
}

type grabOptions struct {
	*commonOptions
}

func parseGrabFlags(fs *flag.FlagSet, args []string) (*grabOptions, error) {
	attachCommonFlags(fs)

	defaultUsage := fs.Usage
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "`%s` %s\n\n", grabCmdName, grabCmdDescription)
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

	opts := &grabOptions{
		commonOptions: commonOpts,
	}
	return opts, nil
}

func grabTabs(opts *grabOptions) error {
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
			return fmt.Errorf("%s\n%v\n", stderr.String(), err)
		}
	}

	// Setup output writer
	var out io.Writer = os.Stdout
	if opts.clipboard {
		out = &clipboard{}
	}
	writer := bufio.NewWriter(out)

	tabs := strings.Split(stdout.String(), "\n")
	for _, tab := range tabs {
		if tab != "" {
			_, err := writer.WriteString(fmt.Sprintf("%s%s\n", opts.prefix, tab))
			if err != nil {
				return fmt.Errorf("failed to write output to buffer: %w", err)
			}
		}
	}
	err := writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}

// Error code indicating tab index out of range
var errCodeEndOfTabs = []byte("(-1719)\n")

func isErrEndOfTabs(buf *bytes.Buffer) bool {
	return bytes.HasSuffix(buf.Bytes(), errCodeEndOfTabs)
}
