package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runTabCmd(cmd *flag.FlagSet, args []string) error {
	opts, err := parseTabFlags(cmd, args)
	if err != nil {
		return err
	}

	err = openTabs(opts)
	if err != nil {
		return err
	}

	return nil
}

type tabOptions struct {
	*commonOptions
	urlReader            io.Reader
	browserArgs          string
	disablePrefixWarning bool
}

func parseTabFlags(fs *flag.FlagSet, args []string) (*tabOptions, error) {
	var (
		urlList = fs.String(
			"urls",
			"",
			fmt.Sprintf("newline-delimited list of URLs, typically the output from the %s command, ignored if -clipboard flag is used", grabCmdName),
		)
		urlFile = fs.String(
			"file",
			"",
			"file containing newline-delimited list of URLs, ignored if -urls or -clipboard flag is used",
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
		fmt.Fprintf(os.Stderr, "`%s` %s\n\n", tabCmdName, tabCmdDescription)
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

	var urlReader io.Reader
	switch {
	case commonOpts.clipboard:
		urlReader = &clipboard{}
	case *urlList != "":
		urlReader = bufio.NewReader(strings.NewReader(*urlList))
	case *urlFile != "":
		f, err := os.Open(*urlFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
		urlReader = bufio.NewReader(f)
	default:
		return nil, errors.New("newline-delimited list of URLs required as a flag argument, from the clipboard, or from a file")
	}

	opts := &tabOptions{
		commonOptions:        commonOpts,
		urlReader:            urlReader,
		browserArgs:          *browserArgs,
		disablePrefixWarning: *disablePrefixWarning,
	}
	return opts, nil
}

func openTabs(opts *tabOptions) error {
	urls, prefixes, err := readURLs(opts.urlReader, func(url string) string {
		return cleanURL(url, opts.prefix)
	})
	if err != nil {
		return fmt.Errorf("failed to read URLs: %w", err)
	}

	if !opts.disablePrefixWarning && warnMismatchingPrefixes(prefixes, opts.prefix) {
		return errUserAbort
	}

	if len(urls) == 0 {
		return errors.New("no URLs provided")
	}

	switch opts.browserApp.name {
	case browserNameChrome, browserNameBrave:
		err := openTabsChromium(opts, urls)
		if err != nil {
			return err
		}
	case browserNameSafari:
		err := openTabsSafari(opts, urls)
		if err != nil {
			return err
		}
	}

	return nil
}

func openTabsChromium(opts *tabOptions, urls []string) error {
	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	newWindowArg := "--new-window" // Open the first URL in a new window
	for _, url := range urls {
		cmd := exec.Command("open", "-na", opts.browserApp.cmdName, "--args", newWindowArg, opts.browserArgs, url)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if opts.verbose {
			log.Printf("executing: %s\n", cmd.String())
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s\n%v\n", stderr.String(), err)
		}
		newWindowArg = "" // Open the remaining URLs as tabs within the window
	}

	return nil
}

func openTabsSafari(opts *tabOptions, urls []string) error {
	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Open a new window
	cmd := exec.Command("open", "-na", opts.browserApp.cmdName, "--args", "--new-window", opts.browserArgs)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if opts.verbose {
		log.Printf("executing: %s\n", cmd.String())
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s\n%v\n", stderr.String(), err)
	}

	// Give Safari a chance to get the new window 1
	time.Sleep(500 * time.Millisecond)

	scriptLayout := "tell application \"" + opts.browserApp.cmdName + "\" to tell window 1 to set URL of %s to \"%s\""
	tabIdx := "tab 1"
	for _, url := range urls {
		cmd := exec.Command("osascript", "-e", fmt.Sprintf(scriptLayout, tabIdx, url))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if opts.verbose {
			log.Printf("executing: %s\n", cmd.String())
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s\n%v\n", stderr.String(), err)
		}
		tabIdx = "(make new tab)"
	}

	// Set last tab as active tab
	lastTabScript := "tell application \"" + opts.browserApp.cmdName + "\" to tell window 1 to set current tab to last tab"
	cmd = exec.Command("osascript", "-e", lastTabScript)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if opts.verbose {
		log.Printf("executing: %s\n", cmd.String())
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s\n%v\n", stderr.String(), err)
	}

	return nil
}

func readURLs(r io.Reader, cleanF func(string) string) ([]string, prefixSet, error) {
	urls := []string{}
	prefixes := newPrefixSet() // Track prefixes to detect potential mismatches

	raw := make([]byte, 1024*1024)
	_, err := r.Read(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read URLs from reader: %w", err)
	}
	rawURLs := string(raw[:])

	for _, url := range strings.Split(rawURLs, "\n") {
		if len(url) == 0 {
			continue
		}

		prefixes.addFrom(url)

		if u := cleanF(url); u != "" {
			urls = append(urls, u)
		}
	}

	return urls, prefixes, nil
}

func cleanURL(url string, prefix string) string {
	return strings.Trim(strings.TrimSpace(strings.TrimPrefix(url, prefix)), "\x00")
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
