package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type restoreOptions struct {
	*commonOptions
	browserArgs string
	urls        []string
}

// TODO: take browser options
func parseRestoreFlags(fs *flag.FlagSet, args []string) (*restoreOptions, error) {
	var (
		urlList     = fs.String("urls", "", "HELP GOES HERE")
		urlFile     = fs.String("file", "", "HELP GOES HERE")
		browserArgs = fs.String("browser-args", "", "HELP GOES HERE")
	)

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	commonOpts, err := parseCommonOptions()
	if err != nil {
		return nil, err
	}

	rawURLs := *urlList
	if rawURLs == "" {
		// Try to parse urls from file
		if *urlFile == "" {
			return nil, errors.New("must provide a newline-delimited list of URLs as a flag argument or with a file")
		}
		raw, err := os.ReadFile(*urlFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
		rawURLs = string(raw[:])
	}

	urls := []string{}
	for _, url := range strings.Split(rawURLs, "\n") {
		if u := strings.TrimPrefix(url, commonOpts.prefix); u != "" {
			urls = append(urls, strings.TrimSpace(u))
		}
	}

	opts := &restoreOptions{
		commonOptions: commonOpts,
		urls:          urls,
		browserArgs:   *browserArgs,
	}
	return opts, nil
}

// TODO: test other browsers (SAFARI)
func restoreTabs(opts *restoreOptions) error {
	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	if len(opts.urls) == 0 {
		return errors.New("no URLs provided")
	}

	newWindowArg := "--new-window" // Open the first URL in a new window
	for _, url := range opts.urls {
		cmd := exec.Command("open", "-na", opts.browserApp.String(), "--args", newWindowArg, opts.browserArgs, url)
		// fmt.Println(cmd.String())
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Error: %s\n%v\n", stderr.String(), err)
		}
		newWindowArg = "" // Open the remaining URLs as tabs within the window
	}

	return nil
}
