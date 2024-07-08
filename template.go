package main

import (
	"io"
	"strings"
	"text/template"
)

const templateURL = "{{.URL}}"

type templatedTabInfoWriter func(*tabInfo) error

func buildTemplateWriteF(w io.Writer, tmpl string) (templatedTabInfoWriter, error) {
	if !strings.HasSuffix(tmpl, "\n") {
		tmpl += "\n"
	}

	t, err := template.New("output").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	writeF := func(tInfo *tabInfo) error {
		return t.Execute(w, tInfo)
	}
	return writeF, nil
}
