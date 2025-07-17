package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

var errEndOfTabs = errors.New("end of tabs")

func execOsaScript(script string, stdout *bytes.Buffer, stderr *bytes.Buffer, verbose bool) error {
	// There is no intention that this implementation be secure so ignore the linter warning
	cmd := exec.Command("osascript", "-e", script) // #nosec
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if verbose {
		log.Printf("executing: %s\n", cmd.String())
	}

	if err := cmd.Run(); err != nil {
		// Check stderr for clean exit on end-of-tabs error
		if stderr.Len() == 0 || isEndOfTabsErrCode(stderr) {
			if verbose {
				log.Print("next tab of window 1 not found, end of tabs")
			}
			return errEndOfTabs
		}
		// Return error if not end-of-tabs
		return fmt.Errorf("%s\n%v\n", stderr.String(), err)
	}

	return nil
}
