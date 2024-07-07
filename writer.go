package main

import (
	"bytes"
	"io"
)

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
