package main

const (
	// Browser names
	browserNameBrave  = "brave"
	browserNameChrome = "chrome"
	browserNameSafari = "safari"
)

// Browser applications by name
var browserApplications = map[string]*browserApplication{
	browserNameBrave: {
		name:    browserNameBrave,
		cmdName: "Brave Browser",
	},
	browserNameChrome: {
		name:    browserNameChrome,
		cmdName: "Google Chrome",
	},
	browserNameSafari: {
		name:    browserNameSafari,
		cmdName: "Safari",
	},
}

type browserApplication struct {
	name    string
	cmdName string
}
