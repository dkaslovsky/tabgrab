package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
	urlWriter io.WriteCloser
}

func parseGrabFlags(fs *flag.FlagSet, args []string) (*grabOptions, error) {
	attachCommonFlags(fs)

	urlFile := fs.String(
		"file",
		"",
		"path for output file containing newline-delimited list of URLs, ignored if -clipboard flag is used",
	)
	quiet := fs.Bool(
		"quiet",
		false,
		"disable console output",
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

	builder := urlWriteCloseRemoverBuilder{}
	if commonOpts.clipboard {
		builder.add(&urlWriteCloseRemover{
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
		builder.add(&urlWriteCloseRemover{
			Writer:  f,
			Closer:  f.Close,
			Remover: func() error { return os.Remove(f.Name()) },
		})
	}
	if !*quiet {
		builder.add(&urlWriteCloseRemover{
			Writer:  os.Stdout,
			Closer:  func() error { return nil },
			Remover: func() error { return nil },
		})
	}
	urlWriter := builder.build()

	opts := &grabOptions{
		commonOptions: commonOpts,
		urlWriter:     urlWriter,
	}
	return opts, nil
}

func grabTabs(opts *grabOptions) error {
	// File cleanup on error
	removeFileOnError := true
	defer func() {
		if removeFileOnError {
			if remover, ok := opts.urlWriter.(*urlWriteCloseRemover); ok {
				_ = remover.Remove()
			}
		}
	}()

	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Script to capture URL of tab i
	tabScript := "tell application \"" + opts.browserApp.cmdName + "\" to get {URL} of tab %d of window 1"

	for i := 0; i < opts.maxTabs; i++ {
		iTabScript := fmt.Sprintf(tabScript, i+1)

		// There is no intention that this implementation be secure so ignore the linter warning
		cmd := exec.Command("osascript", "-e", iTabScript) // #nosec
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if opts.verbose {
			log.Printf("executing: %s\n", cmd.String())
		}

		if err := cmd.Run(); err != nil {
			// Check stderr for clean exit on end-of-tabs error
			if stderr.Len() == 0 || isErrEndOfTabs(&stderr) {
				if opts.verbose {
					log.Printf("tab %d of window 1 not found, end of tabs", i+1)
				}
				break
			}
			// Return error if not end-of-tabs
			return fmt.Errorf("%s\n%v\n", stderr.String(), err)
		}
	}

	writer := bufio.NewWriter(opts.urlWriter)
	defer opts.urlWriter.Close()

	tabs := strings.Split(stdout.String(), "\n")
	for _, tab := range tabs {
		if tab != "" {
			_, err := writer.WriteString(fmt.Sprintf("%s%s\n", opts.prefix, tab))
			if err != nil {
				return fmt.Errorf("failed to write output to buffer: %w", err)
			}
		}
	}
	err := writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	removeFileOnError = false
	return nil
}

// Error code indicating tab index out of range
var errCodeEndOfTabs = []byte("(-1719)\n")

func isErrEndOfTabs(buf *bytes.Buffer) bool {
	return bytes.HasSuffix(buf.Bytes(), errCodeEndOfTabs)
}

type urlWriteCloseRemoverBuilder struct {
	elems []*urlWriteCloseRemover
}

func (builder *urlWriteCloseRemoverBuilder) add(u *urlWriteCloseRemover) {
	builder.elems = append(builder.elems, u)
}

func (builder *urlWriteCloseRemoverBuilder) build() *urlWriteCloseRemover {
	writers := make([]io.Writer, 0, len(builder.elems))
	closers := make([]func() error, 0, len(builder.elems))
	removers := make([]func() error, 0, len(builder.elems))

	for _, elem := range builder.elems {
		writers = append(writers, elem.Writer)
		closers = append(closers, elem.Closer)
		removers = append(removers, elem.Remover)
	}

	u := urlWriteCloseRemover{}
	u.Writer = io.MultiWriter(writers...)
	u.Closer = func() error {
		for _, closer := range closers {
			if err := closer(); err != nil {
				return err
			}
		}
		return nil
	}
	u.Remover = func() error {
		for _, remover := range removers {
			if err := remover(); err != nil {
				return err
			}
		}
		return nil
	}

	return &u
}

type urlWriteCloseRemover struct {
	io.Writer
	Closer  func() error
	Remover func() error
}

func (u *urlWriteCloseRemover) Close() error {
	return u.Closer()
}

func (u *urlWriteCloseRemover) Remove() error {
	return u.Remover()
}
