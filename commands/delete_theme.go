package commands

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/srgrn/postmanerator/configuration"
)

type DeleteTheme struct {
	Config *configuration.Configuration `inject:""`
	Themes interface {
		Delete(theme string) error
	} `inject:""`
}

func (c *DeleteTheme) Is(name string) bool {
	return name == CmdThemesDelete
}

func (c *DeleteTheme) Do() error {
	if len(c.Config.Args) < 3 {
		return errors.New("You must provide the name of the theme you want to delete")
	}

	if err := c.Themes.Delete(c.Config.Args[2]); err != nil {
		return err
	}

	fmt.Fprintln(c.Config.Out, color.GreenString("Theme successfully deleted"))
	return nil
}
