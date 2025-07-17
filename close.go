package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func runCloseCmd(cmd *flag.FlagSet, args []string) error {
	opts, err := parseCloseFlags(cmd, args)
	if err != nil {
		return err
	}

	err = closeTabs(opts)
	if err != nil {
		return err
	}

	return nil
}

type closeOptions struct {
	*commonOptions
	matchVals    []string
	nonMatchVals []string
}

func parseCloseFlags(fs *flag.FlagSet, args []string) (*closeOptions, error) {
	attachCommonFlags(fs)

	var (
		match = fs.String(
			"match",
			"",
			"space delimited list of strings for matching tab URLs to close",
		)
		nonMatch = fs.String(
			"no-match",
			"",
			"space delimited list of strings for non-matching tab URLs to close",
		)
	)

	defaultUsage := fs.Usage
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "`%s` %s\n\n", closeCmdName, closeCmdDescription)
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

	opts := &closeOptions{
		commonOptions: commonOpts,
		matchVals:     strings.Split(*match, " "),
		nonMatchVals:  strings.Split(*nonMatch, " "),
	}
	return opts, nil
}

func closeTabs(opts *closeOptions) error {
	buf := &bytes.Buffer{}
	builder := multiWriteCloseRemoverBuilder{}
	builder.add(&writeCloseRemover{
		Writer:  buf,
		Closer:  func() error { return nil },
		Remover: func() error { return nil },
	})

	err := grabTabs(&grabOptions{
		commonOptions: opts.commonOptions,
		urlWriter:     builder.build(),
		template:      templateURL,
	})
	if err != nil {
		return fmt.Errorf("failed to get tabs for matching: %w", err)
	}

	urlBytes, err := io.ReadAll(buf)
	if err != nil {
		return fmt.Errorf("failed to read tab URLs from buffer: %w", err)
	}

	urls := strings.Split(
		strings.TrimSuffix(string(urlBytes), "\n"), "\n",
	)

	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Script to close URL of tab i
	tabScript := "tell application \"" + opts.browserApp.cmdName + "\" to close tab %d of window 1"

	for i, url := range urls {
		if match(url, opts.matchVals, opts.nonMatchVals) {
			err := execOsaScript(fmt.Sprintf(tabScript, i+1), &stdout, &stderr, opts.verbose)
			if err != nil {
				if errors.Is(err, errEndOfTabs) {
					break
				}
				return fmt.Errorf("failed to close tab: %w", err)
			}
		}
	}

	return nil
}

func match(url string, matchVals []string, nonMatchVals []string) bool {
	numEvals := 0 // Track number of evaluations to avoid returning true on all empty values

	for _, val := range matchVals {
		if len(val) == 0 {
			continue
		}
		numEvals++
		if !strings.Contains(url, val) {
			return false
		}
	}

	for _, val := range nonMatchVals {
		if len(val) == 0 {
			continue
		}
		numEvals++
		if strings.Contains(url, val) {
			return false
		}
	}

	return numEvals > 0
}
