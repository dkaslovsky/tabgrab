package main

const (
	// Browser names
	browserNameAtlas  = "atlas"
	browserNameBrave  = "brave"
	browserNameChrome = "chrome"
	browserNameComet  = "comet"
	browserNameSafari = "safari"
)

// Browser applications by name
var browserApplications = map[string]*browserApplication{
	browserNameAtlas: {
		name:    browserNameAtlas,
		cmdName: "ChatGPT Atlas",
	},
	browserNameBrave: {
		name:    browserNameBrave,
		cmdName: "Brave Browser",
	},
	browserNameChrome: {
		name:    browserNameChrome,
		cmdName: "Google Chrome",
	},
	browserNameComet: {
		name:    browserNameComet,
		cmdName: "Comet",
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
