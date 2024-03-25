# tabgrab
A small (and likely insecure) tool to extract the URL of every open tab in the current browser window on macOS.

`tabgrab` is a macOS-specific commandline tool that prints the URL of all open tabs of the current browser window.

Currently the following browsers are supported:
* Brave (default)
* Chrome
* Safari


Usage:
```bash
$ tabgrab -h
  -browser string
    	browser name (default "brave")
  -max int
    	maximum number of tabs (default 100)
  -prefix string
    	prefix to attach to each tab
```

Example: extract and print all open tabs from the current Chrome window
```bash
$ tabgrab -browser chrome -max 10 -prefix "- "
- https://github.com/dkaslovsky/tabgrab/tree/main
- https://www.espn.com/
- https://news.ycombinator.com/
```
