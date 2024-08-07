package main

import (
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

func runTabsCmd(cmd *flag.FlagSet, args []string) error {
	opts, err := parseTabsFlags(cmd, args)
	if err != nil {
		return err
	}

	err = openTabs(opts)
	if err != nil {
		return err
	}

	return nil
}

type tabsOptions struct {
	*commonOptions
	urlReader            io.ReadCloser
	browserArgs          string
	disablePrefixWarning bool
}

func parseTabsFlags(fs *flag.FlagSet, args []string) (*tabsOptions, error) {
	attachCommonFlags(fs)

	var (
		urlList = fs.String(
			"urls",
			"",
			fmt.Sprintf("newline-delimited list of URLs, typically the output from the %s command, ignored if -clipboard flag is used", grabCmdName),
		)
		urlFile = fs.String(
			"file",
			"",
			"path to file containing newline-delimited list of URLs, ignored if -urls or -clipboard flag is used",
		)
		browserArgs = fs.String(
			"browser-args",
			setStringFlagDefault("", envVarBrowserArgs),
			"optional space-delimited arguments to be passed to the browser",
		)
		disablePrefixWarning = fs.Bool(
			"disable-prefix-warning",
			false,
			"disables warning for potentially mismatched prefix flag and URL prefixes (default false)",
		)
	)

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

	var urlReader *urlReadCloser
	switch {
	case commonOpts.clipboard:
		urlReader = &urlReadCloser{
			Reader: &clipboard{},
			Closer: func() error { return nil },
		}
	case *urlList != "":
		urlReader = &urlReadCloser{
			Reader: strings.NewReader(*urlList),
			Closer: func() error { return nil },
		}
	case *urlFile != "":
		f, err := os.Open(*urlFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read from file: %w", err)
		}
		urlReader = &urlReadCloser{
			Reader: f,
			Closer: f.Close,
		}
	default:
		return nil, errors.New("newline-delimited list of URLs required as a flag argument, from the clipboard, or from a file")
	}

	opts := &tabsOptions{
		commonOptions:        commonOpts,
		urlReader:            urlReader,
		browserArgs:          *browserArgs,
		disablePrefixWarning: *disablePrefixWarning,
	}
	return opts, nil
}

type urlReadCloser struct {
	io.Reader
	Closer func() error
}

func (u *urlReadCloser) Close() error {
	return u.Closer()
}

func openTabs(opts *tabsOptions) error {
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
	default:
		return fmt.Errorf("unrecognized browser: %s", opts.browserApp.name)
	}

	return nil
}

func openTabsChromium(opts *tabsOptions, urls []string) error {
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

func openTabsSafari(opts *tabsOptions, urls []string) error {
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

func readURLs(r io.ReadCloser, cleanF func(string) string) ([]string, prefixSet, error) {
	urls := []string{}
	prefixes := newPrefixSet() // Track prefixes to detect potential mismatches

	raw := make([]byte, 1024*1024)
	_, err := r.Read(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read URLs from reader: %w", err)
	}
	defer r.Close()
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
