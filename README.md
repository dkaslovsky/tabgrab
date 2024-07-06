# tabgrab
A small (and likely insecure) tool to extract or restore the URL of every open tab in the current browser window on macOS.

### Overview and Usage
`tabgrab` is a macOS-specific commandline tool to:
* output the URL of all open tabs of the current browser window (`tabgrab grab`)
* reopen tabs in a new browser window from a list of URLs (`tabgrab tab`)

```
$ tabgrab -h
tabgrab: extract and restore URL tabs to and from the active browser window

Usage:
  grab:		extracts the URL from each tab of the active browser window
  tab:		opens the provided URLs as tabs in a new browser window
  version:	displays application version information

Run `tabgrab <subcommand> -help` for subcommand usage and flags
```

Extract URLs from open tabs with the `grab` command:
```
$ tabgrab grab -h
`grab` extracts the URL from each tab of the active browser window

Usage of grab:
  -browser string
    	browser name (default "chrome")
  -clipboard
    	use clipboard for input/output
  -max int
    	maximum number of tabs (default 100)
  -prefix string
    	optional prefix for each URL
  -verbose
    	enable verbose output
```

Restore tabs from a list of URLs:
```
$ tabgrab tab -h
`tab` opens the provided URLs as tabs in a new browser window

Usage of tab:
  -browser string
    	browser name (default "chrome")
  -browser-args string
    	optional space-delimited arguments to be passed to the browser
  -clipboard
    	use clipboard for input/output
  -disable-prefix-warning
    	disables warning for potentially mismatched prefix flag and URL prefixes (default false)
  -file string
    	file containing newline-delimited list of URLs, ignored if -urls or -clipboard flag is used
  -max int
    	maximum number of tabs (default 100)
  -prefix string
    	optional prefix for each URL
  -urls string
    	newline-delimited list of URLs, typically the output from the grab command, ignored if -clipboard flag is used
  -verbose
    	enable verbose output
```


### Examples

Extract all open tabs from the browser's current window (defaults)
```bash
$ tabgrab grab
https://github.com/dkaslovsky/tabgrab/tree/main
https://www.espn.com/
https://news.ycombinator.com/
```

Extract at most 2 open tabs from the current Safari window and output with a specified prefix
```bash
$ tabgrab grab -browser safari -max 2 -prefix "- "
- https://github.com/dkaslovsky/tabgrab/tree/main
- https://www.espn.com/
```

Clipboard workflows are fully supported. For example, to extract all open tabs to the clipboard:
```bash
$ tabgrab grab -clipboard
```
URL tabs can then be restored from the clipboard:
```bash
$ tabgrab tab -clipboard
```

### Support status for common browsers
* Chrome  - supported (default)
* Brave   - supported
* Safari  - supported
* Firefox - not supported due to compatibility issues with the method used for extracting tab URLs


### Installation
Download the binary for your architecture from the [releases page](https://github.com/dkaslovsky/tabgrab/releases/latest) or clone this repository and build from source.
