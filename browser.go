package main

const (
	// Browser names
	browserNameBrave  = "brave"
	browserNameChrome = "chrome"
	browserNameSafari = "safari"
)

// Browser applications by name
var browserApplications = map[string]*browserApplication{
	browserNameBrave:  {"Brave Browser"},
	browserNameChrome: {"Google Chrome"},
	browserNameSafari: {"Safari"},
}

type browserApplication struct {
	appName string
}

func (b browserApplication) String() string {
	return b.appName
}
