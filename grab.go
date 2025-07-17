package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
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
	urlWriter *writeCloseRemover
	template  string
}

func parseGrabFlags(fs *flag.FlagSet, args []string) (*grabOptions, error) {
	attachCommonFlags(fs)

	var (
		urlFile = fs.String(
			"file",
			"",
			"path for output file containing newline-delimited list of URLs",
		)
		quiet = fs.Bool(
			"quiet",
			false,
			"disable console output",
		)
		template = fs.String(
			"template",
			setStringFlagDefault(defaultTemplate, envVarTemplate),
			"output format specifying tab URL with {{.URL}} tab name with {{.Name}}",
		)
	)

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

	builder := multiWriteCloseRemoverBuilder{}
	if commonOpts.clipboard {
		builder.add(&writeCloseRemover{
			Writer:  &clipboard{},
			Closer:  func() error { return nil },
			Remover: func() error { return nil },
		})
	}
	if *urlFile != "" {
		f, err := os.Create(*urlFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		builder.add(&writeCloseRemover{
			Writer:  f,
			Closer:  f.Close,
			Remover: func() error { return os.Remove(f.Name()) },
		})
	}
	if !*quiet {
		builder.add(&writeCloseRemover{
			Writer:  os.Stdout,
			Closer:  func() error { return nil },
			Remover: func() error { return nil },
		})
	}
	urlWriter := builder.build()

	opts := &grabOptions{
		commonOptions: commonOpts,
		urlWriter:     urlWriter,
		template:      *template,
	}
	return opts, nil
}

func grabTabs(opts *grabOptions) error {
	// Cleanup if exit due to error
	writeCleanup := true
	defer func() {
		if writeCleanup {
			_ = opts.urlWriter.Remove()
		}
	}()

	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Script to capture URL of tab i
	tabScript := "tell application \"" + opts.browserApp.cmdName + "\" to get {URL, NAME} of tab %d of window 1"

	for i := 0; i < opts.maxTabs; i++ {
		err := execOsaScript(fmt.Sprintf(tabScript, i+1), &stdout, &stderr, opts.verbose)
		// TODO: better error handling
		if err == errEndOfTabs {
			fmt.Println("breaking")
			break
		}
		if err != nil {
			fmt.Println("err", err)
			return err
		}
	}

	// Write output
	defer opts.urlWriter.Close()
	writer := bufio.NewWriter(opts.urlWriter)
	writeF, err := buildTemplateWriteF(writer, fmt.Sprintf("%s%s", opts.prefix, opts.template))
	if err != nil {
		return fmt.Errorf("failed to construct writer template: %w", err)
	}
	for _, tab := range parseTabInfo(stdout) {
		if tab.URL != "" || tab.Name != "" {
			err := writeF(tab)
			if err != nil {
				return fmt.Errorf("failed to write output to buffer: %w", err)
			}
		}
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	writeCleanup = false
	return nil
}

type tabInfo struct {
	URL  string
	Name string
}

func parseTabInfo(raw bytes.Buffer) []*tabInfo {
	tInfo := []*tabInfo{}
	tabs := strings.Split(raw.String(), "\n")
	for _, tab := range tabs {
		if tab == "" {
			continue
		}
		info := strings.SplitN(tab, ",", 2)
		tabName := ""
		if len(info) == 2 {
			tabName = strings.TrimSpace(info[1])
		}
		tInfo = append(tInfo, &tabInfo{
			URL:  info[0],
			Name: tabName,
		})
	}
	return tInfo
}
