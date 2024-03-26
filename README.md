# tabgrab
A small (and likely insecure) tool to extract the URL of every open tab in the current browser window on macOS.

`tabgrab` is a macOS-specific commandline tool that prints the URL of all open tabs of the current browser window.

Support status for common browsers:
* Brave - supported (default)
* Chrome - supported
* Safari - supported
* Firefox - not supported due to compatibility issues with the method used for extracting tab URLs


### Usage
```bash
$ tabgrab -h
  -browser string
    	browser name (default "brave")
  -max int
    	maximum number of tabs (default 100)
  -prefix string
    	optional prefix to attach to each URL
  -version
    	display version
```

### Examples

Extract all open tabs from the browser's current window (defaults)
```bash
$ tabgrab
https://github.com/dkaslovsky/tabgrab/tree/main
https://www.espn.com/
https://news.ycombinator.com/
```

Extract at most 2 open tabs from the current Chrome window with a prefix
```bash
$ tabgrab -browser chrome -max 2 -prefix "- "
- https://github.com/dkaslovsky/tabgrab/tree/main
- https://www.espn.com/
```
