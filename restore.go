package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type restoreOptions struct {
	*commonOptions
	urls []string
}

func parseRestoreFlags(fs *flag.FlagSet, args []string) (*restoreOptions, error) {
	urlList := fs.String("urls", "", "HELP GOES HERE")
	urlFile := fs.String("file", "", "HELP GOES HERE")

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
			return nil, errors.New("must provide a newline-delimited list of urls as a flag argument or with a file")
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
	}
	return opts, nil
}

func restoreTabs(opts *restoreOptions) error {
	fmt.Printf("\n%+v\n", opts.commonOptions)
	fmt.Printf("\n%+v\n", opts.urls)
	return nil
}
