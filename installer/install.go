package installer

import "github.com/deviceio/shared/install"

func Install(org string, huburl string, hubport int, hubSelfSigned bool) error {
	i := &install.Installer{
		Title: "Installer",
		Steps: []install.Step{
			&step{},
		},
	}

	return i.RunInstall(install.Args{
		"org":                 org,
		"transportHost":       huburl,
		"transportPort":       hubport,
		"transportSelfSigned": hubSelfSigned,
	})
}
