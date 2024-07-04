package main

import (
	"flag"
	"fmt"
	"strings"
)

type openOptions struct {
	*commonOptions
	urls []string
}

func parseOpenFlags(fs *flag.FlagSet, args []string) (*openOptions, error) {
	urls := fs.String("urls", "", "HELP GOES HERE")
	file := fs.String("file", "", "HELP GOES HERE")
	_ = file

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	commonOpts, err := setCommonOptions()
	if err != nil {
		return nil, err
	}

	openURLs := []string{}
	if *urls != "" {
		for _, u := range strings.Split(*urls, "\n") {
			openURLs = append(openURLs, strings.TrimPrefix(u, commonOpts.prefix))
		}
	}

	opts := &openOptions{
		commonOptions: commonOpts,
		urls:          openURLs,
	}
	return opts, nil
}

func openTabs(opts *openOptions) error {
	fmt.Printf("===\n%+v\n", opts)

	return nil
}
