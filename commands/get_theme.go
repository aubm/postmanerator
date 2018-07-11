package commands

import (
	"errors"
	"fmt"

	"github.com/aubm/postmanerator/configuration"
	"github.com/fatih/color"
)

type GetTheme struct {
	Config *configuration.Configuration `inject:""`
	Themes interface {
		Download(themeName string) error
	} `inject:""`
}

func (c *GetTheme) Is(name string) bool {
	return name == CmdThemesGet
}

func (c *GetTheme) Do() error {
	if len(c.Config.Args) < 3 {
		return errors.New("You must provide the name or the URL of the theme you want to download")
	}

	if err := c.Themes.Download(c.Config.Args[2]); err != nil {
		return err
	}

	fmt.Fprintln(c.Config.Out, color.GreenString("Theme successfully downloaded"))
	return nil
}
