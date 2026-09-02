package config

import (
	"errors"
)

func LoadArgs() CLIArgs {
	// first check if we have a config ready to use,
	// if not get cli args,
	// if still missing required - prompt for configuration
	configArgs, err := LoadConfig()
	cliArgs := ParseCliArgs()

	// either the user needs to configure or the user wants to configure
	if (err != nil && errors.Is(err, NotConfiguredErr)) || cliArgs.Configure {
		HandleConfigCreation(&configArgs)
	}
	if cliArgs.ModelId == "" {
		cliArgs.ModelId = configArgs.ModelId
	}

	if cliArgs.Provider == "" {
		cliArgs.Provider = configArgs.Provider
	}
	return cliArgs
}
