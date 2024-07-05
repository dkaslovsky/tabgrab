package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type restoreOptions struct {
	*commonOptions
	browserArgs string
	urls        []string
}

func parseRestoreFlags(fs *flag.FlagSet, args []string) (*restoreOptions, error) {
	var (
		urlList = fs.String(
			"urls",
			"",
			fmt.Sprintf("newline-delimited list of URLs, typically the output from the `%s` command", saveCmdName),
		)
		urlFile = fs.String(
			"file",
			"",
			"file containing newline-delimited list of URLs, ignored if --urls flag is used",
		)
		browserArgs = fs.String(
			"browser-args",
			"",
			"optional space-delimited arguments to be passed to the browser",
		)
		disablePrefixWarning = fs.Bool(
			"disable-prefix-warning",
			false,
			"disables warning for potentially mismatched prefix flag and URL prefixes (default false)",
		)
	)

	attachCommonFlags(fs)

	defaultUsage := fs.Usage
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "`%s` opens the provided URLs as tabs in a new browser window\n\n", restoreCmdName)
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

	// Temporary error for unsupported browser
	if commonOpts.browserApp.name == browserNameSafari {
		return nil, fmt.Errorf("%s is not yet supported for this subcommand", browserNameSafari)
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
	prefixes := newPrefixSet() // Track prefixes to detect potential mismatches
	for _, url := range strings.Split(rawURLs, "\n") {
		if len(url) == 0 {
			continue
		}

		prefixes.addFrom(url)

		if u := strings.TrimSpace(strings.TrimPrefix(url, commonOpts.prefix)); u != "" {
			urls = append(urls, u)
		}
	}

	if !(*disablePrefixWarning) && warnMismatchingPrefixes(prefixes, commonOpts.prefix) {
		return nil, errUserAbort
	}

	opts := &restoreOptions{
		commonOptions: commonOpts,
		urls:          urls,
		browserArgs:   *browserArgs,
	}
	return opts, nil
}

func restoreTabs(opts *restoreOptions) error {
	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	if len(opts.urls) == 0 {
		return errors.New("no URLs provided")
	}

	newWindowArg := "--new-window" // Open the first URL in a new window
	for _, url := range opts.urls {
		cmd := exec.Command("open", "-na", opts.browserApp.cmdName, "--args", newWindowArg, opts.browserArgs, url)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if opts.verbose {
			log.Printf("executing: %s\n", cmd.String())
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Error: %s\n%v\n", stderr.String(), err)
		}

		newWindowArg = "" // Open the remaining URLs as tabs within the window
	}

	return nil
}

var errUserAbort = errors.New("user aborted")

func warnMismatchingPrefixes(prefixes prefixSet, targetPrefix string) bool {
	// Do not warn and prompt for abort if no mismatch exists
	if !checkPrefixMismatch(prefixes, targetPrefix) {
		return false
	}

	foundPrefix, ok := prefixes.pop()
	if !ok {
		return false // Do not abort if a prefix cannot be popped from the set
	}

	if targetPrefix == "" {
		fmt.Printf("Warning: all URLs contain prefix \"%s\" but prefix flag not provided\n", foundPrefix)
	} else {
		fmt.Printf("Warning: all URLs contain prefix \"%s\" but prefix flag \"%s\" does not match\n", foundPrefix, targetPrefix)
	}

	var userInput string
	for userInput != "Y" && userInput != "n" {
		fmt.Print("Continue? [Y/n]: ")
		fmt.Scanln(&userInput)
	}

	return userInput == "n"
}
