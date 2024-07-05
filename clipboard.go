package main

import (
	"os/exec"
)

// Credit: modified from https://stackoverflow.com/questions/73812535/how-to-get-copied-text-from-clipboard-on-golang-mac

// clipboard implements the io.ReadWriter interface
type clipboard struct{}

func (c *clipboard) Write(content []byte) (int, error) {
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	if _, err := in.Write(content); err != nil {
		return 0, err
	}
	if err := in.Close(); err != nil {
		return 0, err
	}
	return len(content), cmd.Wait()
}

func (c *clipboard) Read(buf []byte) (int, error) {
	cmd := exec.Command("pbpaste")
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	return copy(buf, out), nil
}
