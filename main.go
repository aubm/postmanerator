package main

import (
	"fmt"
	"os"

	"github.com/aubm/postmanerator/commands"
	"github.com/aubm/postmanerator/configuration"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
)

const (
	documentationUrl = "https://github.com/aubm/postmanerator"
	cmdDefault       = "cmd_default"
	cmdThemesList    = "cmd_themes_list"
	cmdThemesGet     = "cmd_themes_get"
	cmdThemesDelete  = "cmd_themes_delete"
	cmdUnknown       = "cmd_unknown"
)

var (
	config             = configuration.Config
	errUnknownCmd      = fmt.Errorf(`Command not found, please see the documentation at %s`, documentationUrl)
	defaultCommand     = commands.Default{}
	getThemeCommand    = commands.GetTheme{}
	deleteThemeCommand = commands.DeleteTheme{}
	listThemesCommand  = commands.ListThemes{}
)

func init() {
	checkAndPrintErr(_init())
}

func _init() error {
	if err := inject.Populate(&config, &defaultCommand, &getThemeCommand, &deleteThemeCommand); err != nil {
		return fmt.Errorf("app initialization failed: %v", err)
	}
	return nil
}

func main() {
	checkAndPrintErr(configuration.InitErr)
	checkAndPrintErr(_main())
}

func _main() (err error) {
	switch evaluateUserCommand() {
	case cmdDefault:
		err = defaultCommand.Do()
	case cmdThemesGet:
		err = getThemeCommand.Do()
	case cmdThemesDelete:
		err = deleteThemeCommand.Do()
	case cmdThemesList:
		err = listThemesCommand.Do()
	default:
		err = errUnknownCmd
	}
	return
}

func evaluateUserCommand() string {
	if len(config.Args) == 0 {
		return cmdDefault
	}

	switch config.Args[0] {
	case "themes":
		if len(config.Args) < 2 {
			return cmdThemesList
		}
		switch config.Args[1] {
		case "get":
			return cmdThemesGet
		case "delete":
			return cmdThemesDelete
		case "list":
			return cmdThemesList
		}
	}

	return cmdUnknown
}

func checkAndPrintErr(err error) {
	if err == nil {
		return
	}

	fmt.Println(color.RedString(err.Error()))
	os.Exit(1)
}
