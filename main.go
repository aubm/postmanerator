package main

import (
	"fmt"
	"os"

	"github.com/aubm/postmanerator/commands"
	"github.com/aubm/postmanerator/configuration"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
)

var (
	config             = configuration.Config
	errUnknownCmd      = fmt.Errorf("Command not found, please see the documentation at https://github.com/aubm/postmanerator")
	defaultCommand     = &commands.Default{}
	getThemeCommand    = &commands.GetTheme{}
	deleteThemeCommand = &commands.DeleteTheme{}
	listThemesCommand  = &commands.ListThemes{}
	availableCommands  = []commands.Command{}
)

func init() {
	checkAndPrintErr(_init())
}

func _init() error {
	if err := inject.Populate(&config, defaultCommand, getThemeCommand, deleteThemeCommand, listThemesCommand); err != nil {
		return fmt.Errorf("app initialization failed: %v", err)
	}
	availableCommands = append(availableCommands,
		defaultCommand,
		getThemeCommand,
		deleteThemeCommand,
		listThemesCommand,
	)
	return nil
}

func main() {
	checkAndPrintErr(configuration.InitErr)
	checkAndPrintErr(_main())
}

func _main() (err error) {
	userCommand := evaluateUserCommand()
	for _, availableCommand := range availableCommands {
		if availableCommand.CanHandle(userCommand) {
			return availableCommand.Do()
		}
	}
	return errUnknownCmd
}

func evaluateUserCommand() string {
	if len(config.Args) == 0 {
		return commands.CmdDefault
	}

	switch config.Args[0] {
	case "themes":
		if len(config.Args) < 2 {
			return commands.CmdThemesList
		}
		switch config.Args[1] {
		case "get":
			return commands.CmdThemesGet
		case "delete":
			return commands.CmdThemesDelete
		case "list":
			return commands.CmdThemesList
		}
	}

	return commands.CmdUnknown
}

func checkAndPrintErr(err error) {
	if err == nil {
		return
	}

	fmt.Println(color.RedString(err.Error()))
	os.Exit(1)
}
