package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Cap the number of tabs captured at a reasonable default
const maxTabs = 100

func main() {
	// Define a buffer to capture stdout
	var stdout bytes.Buffer

	// Script to capture URL of tab i
	tabScript := "tell application \"Brave Browser\" to get {URL} of tab %d of window 1"

	for i := 1; i < maxTabs; i++ {

		iTabScript := fmt.Sprintf(tabScript, i)

		// There is no intention that this implementation be secure so ignore the linter warning
		cmd := exec.Command("osascript", "-e", iTabScript) // #nosec
		cmd.Stdout = &stdout

		if err := cmd.Run(); err != nil {
			break
		}

	}

	fmt.Printf("%s", stdout.String())
}
