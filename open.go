package main

import (
	"flag"
	"fmt"
	"strings"
)

type restoreOptions struct {
	*commonOptions
	urls []string
}

func parseRestoreFlags(fs *flag.FlagSet, args []string) (*restoreOptions, error) {
	urls := fs.String("urls", "", "HELP GOES HERE")
	file := fs.String("file", "", "HELP GOES HERE")
	_ = file

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	commonOpts, err := parseCommonOptions()
	if err != nil {
		return nil, err
	}

	restoreURLs := []string{}
	if *urls != "" {
		for _, url := range strings.Split(*urls, "\n") {
			if u := strings.TrimPrefix(url, commonOpts.prefix); u != "" {
				restoreURLs = append(restoreURLs, u)
			}
		}
	}

	opts := &restoreOptions{
		commonOptions: commonOpts,
		urls:          restoreURLs,
	}
	return opts, nil
}

func restoreTabs(opts *restoreOptions) error {
	fmt.Printf("\n%+v\n", opts)

	return nil
}
