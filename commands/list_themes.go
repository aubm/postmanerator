package commands

import (
	"os"

	"github.com/aubm/postmanerator/configuration"
	"github.com/aubm/postmanerator/theme"
)

type ListThemes struct {
	Config *configuration.Configuration `inject:""`
}

func (c *ListThemes) CanHandle(name string) bool {
	return name == CmdThemesList
}

func (c *ListThemes) Do() error {
	return theme.ListThemes(os.Stdout, c.Config.ThemesDirectory)
}
