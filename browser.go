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
		Name:    browserNameBrave,
		appName: "Brave Browser",
	},
	browserNameChrome: {
		Name:    browserNameChrome,
		appName: "Google Chrome",
	},
	browserNameSafari: {
		Name:    browserNameSafari,
		appName: "Safari",
	},
}

type browserApplication struct {
	Name    string
	appName string
}

func (b browserApplication) String() string {
	return b.appName
}
