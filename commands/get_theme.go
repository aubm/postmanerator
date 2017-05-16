package commands

import (
	"errors"

	"github.com/aubm/postmanerator/configuration"
	"github.com/aubm/postmanerator/theme"
	"github.com/fatih/color"
)

type GetTheme struct {
	Config *configuration.Configuration `inject:""`
}

func (c *GetTheme) Do() error {
	if len(c.Config.Args) < 3 {
		return errors.New("You must provide the name or the URL of the theme you want to download")
	}

	if err := theme.GitClone(c.Config.Args[2], c.Config.LocalName, c.Config.ThemesRepository, theme.DefaultCloner{ThemesDirectory: c.Config.ThemesDirectory}); err != nil {
		return err
	}

	c.Config.Out.Write([]byte(color.GreenString("Theme successfully downloaded")))
	return nil
}
