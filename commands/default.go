package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"text/template"

	"github.com/aubm/postmanerator/configuration"
	"github.com/aubm/postmanerator/postman"
	"github.com/aubm/postmanerator/theme"
	"github.com/aubm/postmanerator/theme/helper"
	"github.com/fatih/color"
	"github.com/howeyc/fsnotify"
)

type Default struct {
	Config *configuration.Configuration `inject:""`
}

func (c *Default) Do() error {
	var (
		err error
		out *os.File = os.Stdout
	)

	if c.Config.CollectionFile == "" {
		return errors.New("You must provide a collection using the -collection flag")
	}

	var env map[string]string
	if c.Config.EnvironmentFile != "" {
		env, err = postman.EnvironmentFromFile(c.Config.EnvironmentFile)
		if err != nil {
			return fmt.Errorf("Failed to open environment file: %v", err)
		}
	}

	if c.Config.OutputFile != "" {
		out, err = os.Create(c.Config.OutputFile)
		if err != nil {
			return fmt.Errorf("Failed to create output: %v", err)
		}
		defer out.Close()
	}

	col, err := postman.CollectionFromFile(c.Config.CollectionFile, postman.CollectionOptions{
		IgnoredRequestHeaders:  postman.HeadersList(c.Config.IgnoredRequestHeaders.Values),
		IgnoredResponseHeaders: postman.HeadersList(c.Config.IgnoredResponseHeaders.Values),
		EnvironmentVariables:   env,
	})
	if err != nil {
		return fmt.Errorf("Failed to parse collection file: %v", err)
	}

	col.ExtractStructuresDefinition()

	themePath, err := c.getThemePath()
	if err != nil {
		return errors.New("The theme was not found")
	}
	themeFiles, err := theme.ListThemeFiles(themePath)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			c.Config.Out.Write([]byte(color.RedString("FAIL. %v\n", r)))
		}
	}()
	c.generate(c.Config.Out, out, themeFiles, col)

	if c.Config.Watch {
		return c.watchDir(c.Config.Out, themePath, func() { c.generate(c.Config.Out, out, themeFiles, col) })
	}

	return nil
}

func (c *Default) getThemePath() (string, error) {
	themePath, err := theme.GetThemePath(c.Config.UsedTheme, c.Config.ThemesDirectory)
	if err != nil {
		if ok, _ := regexp.MatchString(`\/|\\`, c.Config.UsedTheme); ok == false {
			c.Config.Out.Write([]byte(color.BlueString("Theme '%v' not found, trying to download it...\n", c.Config.UsedTheme)))
			if err := theme.GitClone(c.Config.UsedTheme, "", c.Config.ThemesRepository, theme.DefaultCloner{ThemesDirectory: c.Config.ThemesDirectory}); err == nil {
				return theme.GetThemePath(c.Config.UsedTheme, c.Config.ThemesDirectory)
			}
		}
		return "", err
	}
	return themePath, nil
}

func (c *Default) generate(out io.Writer, f *os.File, themeFiles []string, col *postman.Collection) {
	out.Write([]byte("Generating output ... "))
	f.Truncate(0)
	templates := template.Must(template.New("").Funcs(helper.GetFuncMap()).ParseFiles(themeFiles...))
	err := templates.ExecuteTemplate(f, "index.tpl", *col)
	if err != nil {
		out.Write([]byte(color.RedString("FAIL. %v\n", err)))
		return
	}
	out.Write([]byte(color.GreenString("SUCCESS.\n")))
}

func (c *Default) watchDir(out io.Writer, dir string, action func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	defer func() {
		<-done
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				out.Write([]byte(color.RedString("FAIL. %v\n", r)))
			}
			c.watchDir(out, dir, action)
		}()

		for {
			select {
			case ev := <-watcher.Event:
				if !ev.IsAttrib() {
					action()
				}
			case err := <-watcher.Error:
				out.Write([]byte(color.RedString("FAIL. %v\n", err)))
			}
		}
	}()

	if err := watcher.Watch(dir); err != nil {
		return fmt.Errorf("Failed to watch theme directory: %v", err)
	}

	return nil
}
