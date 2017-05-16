package commands

import (
	"errors"

	"github.com/aubm/postmanerator/configuration"
	"github.com/aubm/postmanerator/theme"
	"github.com/fatih/color"
)

type DeleteTheme struct {
	Config *configuration.Configuration `inject:""`
}

func (c *DeleteTheme) Do() error {
	if len(c.Config.Args) < 3 {
		return errors.New("You must provide the name of the theme you want to delete")
	}

	if err := theme.Delete(c.Config.Args[2], c.Config.ThemesDirectory); err != nil {
		return err
	}

	c.Config.Out.Write([]byte(color.GreenString("Theme successfully deleted")))
	return nil
}
