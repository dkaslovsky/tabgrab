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

	var (
		urlFile = fs.String(
			"file",
			"",
			"path for output file containing newline-delimited list of URLs, ignored if -clipboard flag is used",
		)
		quiet = fs.Bool(
			"quiet",
			false,
			"disable console output",
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
	}
	return opts, nil
}

func grabTabs(opts *grabOptions) error {
	// Cleanup if exit due to error
	writeCleanup := true
	defer func() {
		if writeCleanup {
			if remover, ok := opts.urlWriter.(*writeCloseRemover); ok {
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

	writeCleanup = false
	return nil
}

type writeCloseRemover struct {
	io.Writer
	Closer  func() error
	Remover func() error
}

func (w *writeCloseRemover) Close() error {
	return w.Closer()
}

func (w *writeCloseRemover) Remove() error {
	return w.Remover()
}

type multiWriteCloseRemoverBuilder struct {
	elems []*writeCloseRemover
}

func (builder *multiWriteCloseRemoverBuilder) add(w *writeCloseRemover) {
	builder.elems = append(builder.elems, w)
}

func (builder *multiWriteCloseRemoverBuilder) build() *writeCloseRemover {
	writers := make([]io.Writer, 0, len(builder.elems))
	closers := make([]func() error, 0, len(builder.elems))
	removers := make([]func() error, 0, len(builder.elems))

	for _, elem := range builder.elems {
		writers = append(writers, elem.Writer)
		closers = append(closers, elem.Closer)
		removers = append(removers, elem.Remover)
	}

	w := writeCloseRemover{}
	w.Writer = io.MultiWriter(writers...)
	w.Closer = func() error {
		for _, closer := range closers {
			if err := closer(); err != nil {
				return err
			}
		}
		return nil
	}
	w.Remover = func() error {
		for _, remover := range removers {
			if err := remover(); err != nil {
				return err
			}
		}
		return nil
	}
	return &w
}

// Error code indicating tab index out of range
var errCodeEndOfTabs = []byte("(-1719)\n")

func isErrEndOfTabs(buf *bytes.Buffer) bool {
	return bytes.HasSuffix(buf.Bytes(), errCodeEndOfTabs)
}
