# tabgrab
A small command-line tool to extract or restore the URL of every open tab in the current browser window on macOS.

### Overview and Usage
`tabgrab` is a macOS-specific command-line tool to:
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
  -file string
    	path for output file containing newline-delimited list of URLs
  -max int
    	maximum number of tabs (default 100)
  -prefix string
    	optional prefix for each URL
  -quiet
    	disable console output
  -template string
    output format specifying tab URL with {{.URL}} tab name with {{.Name}} (default "{{.URL}}")
  -verbose
    	enable verbose output
```

Restore tabs from a list of URL with the `tab` command:
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
    	path to file containing newline-delimited list of URLs, ignored if -urls or -clipboard flag is used
  -max int
    	maximum number of tabs (default 100)
  -prefix string
    	optional prefix for each URL
  -urls string
    	newline-delimited list of URLs, typically the output from the grab command, ignored if -clipboard flag is used
  -verbose
    	enable verbose output
```
The following environment variables can be used to change default flag values:
* `TABGRAB_BROWSER`: sets the default for the `browser` flag
* `TABGRAB_BROWSER_ARGS`: sets the default for the `browser-args` flag
* `TABGRAB_PREFIX`: sets the the default for the `prefix` flag
* `TABGRAB_TEMPLATE`: sets the the default for the `template` flag

### Examples

#### Using stdout
Extract all open tabs from the browser's current window (defaults)
```
$ tabgrab grab
https://github.com/dkaslovsky/tabgrab/tree/main
https://www.espn.com/
https://news.ycombinator.com/
```

Extract at most 2 open tabs from the current Safari window and output with a specified prefix
```
$ tabgrab grab -browser safari -max 2 -prefix "- "
- https://github.com/dkaslovsky/tabgrab/tree/main
- https://www.espn.com/
``` 

#### Output template
Extract the name and URL from all open tabs and output in markdown format
```
$ tabgrab grab -template "[{{.Name}}]({{.URL}})"
[ESPN - Serving Sports Fans. Anytime. Anywhere.](https://www.espn.com/)
[Hacker News](https://news.ycombinator.com/)
```

Combine the above with a prefix
```
$ tabgrab grab -template "[{{.Name}}]({{.URL}})" -prefix "* "
* [ESPN - Serving Sports Fans. Anytime. Anywhere.](https://www.espn.com/)
* [Hacker News](https://news.ycombinator.com/)
```

A prefix can instead be included in the template string if so desired, but must be specified using the `-prefix` flag when restoring tabs with the `tab` command.


#### Using the clipboard
To extract all open tabs to the clipboard:
```
$ tabgrab grab -quiet -clipboard
```
URL tabs can then be restored from the clipboard:
```
$ tabgrab tab -quiet -clipboard
```

#### Using a file
```
$ tabgrab grab -quiet -file "my-tabs.txt"
```
URL tabs can then be restored from the clipboard:
```
$ tabgrab tab -quiet -file "my-tabs.txt"
```

#### Multiple outputs
Output is written to each of stdout, the clipboard, and a specified file by including both the `-file` and `-clipboard` flags and removing the `-quiet` flag.

### Support status for common browsers
* Chrome  - supported (default)
* Brave   - supported
* Safari  - supported
* Firefox - not supported due to compatibility issues with the method used for extracting tab URLs


### Security
`tabgrab` makes no guarantees about security and executes shell commands/Apple Scripts to open browser applications.

### Installation
Download the binary for your architecture from the [releases page](https://github.com/dkaslovsky/tabgrab/releases/latest) or clone this repository and build from source.
