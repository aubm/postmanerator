package commands

import (
	"fmt"

	"github.com/srgrn/postmanerator/configuration"
)

type ListThemes struct {
	Config *configuration.Configuration `inject:""`
	Themes interface {
		List() ([]string, error)
	} `inject:""`
}

func (c *ListThemes) Is(name string) bool {
	return name == CmdThemesList
}

func (c *ListThemes) Do() error {
	themeList, err := c.Themes.List()
	if err != nil {
		return err
	}

	for _, theme := range themeList {
		fmt.Fprintln(c.Config.Out, theme)
	}

	return nil
}
