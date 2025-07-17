package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
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
	matchSubStr string // single exact match string for now
}

func parseCloseFlags(fs *flag.FlagSet, args []string) (*closeOptions, error) {
	attachCommonFlags(fs)

	match := fs.String(
		"match",
		"",
		"(sub)string to match tab URL",
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
		matchSubStr:   *match,
	}
	return opts, nil
}

func closeTabs(opts *closeOptions) error {
	if opts.matchSubStr == "" {
		return nil
	}

	closeAllTabs := opts.matchSubStr == "*"

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
		// TODO: better error handling
		panic(err)
	}

	data, err := io.ReadAll(buf)
	if err != nil {
		// TODO: better error handling
		panic(err)
	}

	matchBytes := []byte(opts.matchSubStr)
	eolBytes := []byte("\n")

	urls := bytes.Split(bytes.TrimSuffix(data, eolBytes), eolBytes)

	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Script to close URL of tab i
	tabScript := "tell application \"" + opts.browserApp.cmdName + "\" to close tab %d of window 1"

	for i, url := range urls {
		if closeAllTabs || bytes.Contains(url, matchBytes) {
			err := execOsaScript(fmt.Sprintf(tabScript, i+1), &stdout, &stderr, opts.verbose)
			// TODO: better error handling
			if err == errEndOfTabs {
				break
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}
