package command

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/standup-raven/standup-raven/server/config"
	"github.com/standup-raven/standup-raven/server/util"
	"strings"
)

// Master is the driver command for all other commands
// All other slash commands are run as /standup <command-name> [command-args]
func Master() *Config {
	return &Config{
		Command: &model.Command{
			Trigger:          config.CommandPrefix,
			AutoComplete:     true,
			AutoCompleteDesc: "Available commands: " + strings.Join(getAvailableCommands(), ", "),
		},
		HelpText: "",
		Validate: validateCommandMaster,
		Execute:  executeCommandMaster,
	}
}

func getAvailableCommands() []string {
	availableCommands := []string{}
	for command := range commands {
		availableCommands = append(availableCommands, command)
	}
	return availableCommands
}

func validateCommandMaster(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	// validate that a command is specified
	if len(args) == 0 {
		return util.SendEphemeralText("Please specify a command")
	}

	subCommand := args[0]
	subCommandCommand, ok := commands[subCommand]

	// validate sub-command exists
	if !ok {
		return util.SendEphemeralText("Invalid command: " + subCommand)
	}

	// add sub-command in props so we don't need to extract it again
	context.Props["subCommand"] = subCommandCommand
	context.Props["subCommandArgs"] = args[1:]

	// run validation for sub-command
	if response, appErr := subCommandCommand.Validate(args[1:], context); response != nil || appErr != nil {
		return response, appErr
	}

	// all okay
	return nil, nil
}

func executeCommandMaster(args []string, context Context) (*model.CommandResponse, *model.AppError) {
	subCommand := context.Props["subCommand"].(*Config)
	subCommandArgs := context.Props["subCommandArgs"].([]string)
	return subCommand.Execute(subCommandArgs, context)
}
