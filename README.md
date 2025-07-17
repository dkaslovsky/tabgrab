# tabgrab
A small command-line tool to extract or restore the URL of every open tab in the current browser window on macOS.

</br>

### Overview and Usage
`tabgrab` is a macOS-specific command-line tool to:
* output the URL of all open tabs of the current browser window (`tabgrab grab`)
* reopen tabs in a new browser window from a list of URLs (`tabgrab tabs`)

```
$ tabgrab -h
tabgrab: extract and restore URL tabs to and from the active browser window

Usage:
  grab:		extracts the URL from each tab of the active browser window
  tabs:		opens the provided URLs as tabs in a new browser window
  close:	closes tabs based on URL matching
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

Restore tabs from a list of URL with the `tabs` command:
```
$ tabgrab tabs -h
`tabs` opens the provided URLs as tabs in a new browser window

Usage of tabs:
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

Close tabs based on URL matching with the `close` command:
```
$ tabgrab close -h
`close` closes tabs based on URL matching

Usage of close:
  -browser string
    	browser name (default "brave")
  -clipboard
    	use clipboard for input/output
  -match string
    	space delimited list of strings for matching tab URLs to close
  -max int
    	maximum number of tabs (default 100)
  -no-match string
    	space delimited list of strings for non-matching tab URLs to close
  -prefix string
    	optional prefix for each URL
  -verbose
```

The following environment variables can be used to change default flag values:
* `TABGRAB_BROWSER`: sets the default for the `browser` flag
* `TABGRAB_BROWSER_ARGS`: sets the default for the `browser-args` flag
* `TABGRAB_PREFIX`: sets the the default for the `prefix` flag
* `TABGRAB_TEMPLATE`: sets the the default for the `template` flag

</br>

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

A prefix can instead be included in the template string if so desired, but must be specified using the `-prefix` flag when restoring tabs with the `tabs` command.

#### Using the clipboard
To extract all open tabs to the clipboard:
```
$ tabgrab grab -quiet -clipboard
```
URL tabs can then be restored from the clipboard:
```
$ tabgrab tabs -quiet -clipboard
```

#### Using a file
```
$ tabgrab grab -quiet -file "my-tabs.txt"
```
URL tabs can then be restored from the clipboard:
```
$ tabgrab tabs -quiet -file "my-tabs.txt"
```

#### Multiple outputs
Output is written to each of stdout, the clipboard, and a specified file by including both the `-file` and `-clipboard` flags and removing the `-quiet` flag.

</br>

#### Close tabs
Close tabs with URLs containing "foo":
```
$ tabgrab close -match "foo"
```
Close tabs with URLs containing "foo" and "bar":
```
$ tabgrab close -match "foo bar"
```
Close tabs with URLs not containing "foo":
```
$ tabgrab close -no-match "foo"
```
Close tabs with URLs containing "foo" but not containing "bar":
```
$ tabgrab close -match "foo" -no-match "bar"
```


</br>

### Support status for common browsers
* Chrome  - supported (default)
* Brave   - supported
* Safari  - supported
* Firefox - not supported due to compatibility issues with the method used for extracting tab URLs

</br>

### Installation Options
* Download a pre-built binary (see [releases](https://github.com/dkaslovsky/tabgrab/releases/latest)):
  
    ARM:
    ```
    $ curl -o tabgrab -L https://github.com/dkaslovsky/tabgrab/releases/latest/download/tabgrab_darwin_arm64
    ```
    AMD:
    ```
    $ curl -o tabgrab -L https://github.com/dkaslovsky/tabgrab/releases/latest/download/tabgrab_darwin_amd64
    ```

* Install using Go:
    ```
    $ go install github.com/dkaslovsky/tabgrab@latest
    ```

* Build from source by cloning this repository and running `go build` in the `tabgrab` root directory.

</br>

### Security
`tabgrab` makes no guarantees about security and executes shell commands/Apple Scripts to open browser applications.
