package commands

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aubm/postmanerator/configuration"
	"github.com/aubm/postmanerator/postman"
	"github.com/aubm/postmanerator/themes"
	"github.com/fatih/color"
	"github.com/howeyc/fsnotify"
)

type Default struct {
	Config *configuration.Configuration `inject:""`
	Themes interface {
		Open(themeName string) (*themes.Theme, error)
		Download(themeName string) error
	} `inject:""`
	CollectionBuilder interface {
		FromFile(file string, options postman.BuilderOptions) (postman.Collection, error)
	} `inject:""`
	EnvironmentBuilder interface {
		FromFile(file string) (postman.Environment, error)
	} `inject:""`
	Renderer interface {
		Render(w io.Writer, theme *themes.Theme, collection postman.Collection) error
	} `inject:""`
}

func (c *Default) Is(name string) bool {
	return name == CmdDefault
}

func (c *Default) Do() error {
	if err := c.validateUserInput(); err != nil {
		return err
	}

	postmanEnvironment, err := c.getPostmanEnvironment()
	if err != nil {
		return err
	}

	postmanCollection, err := c.getPostmanCollection(postmanEnvironment)
	if err != nil {
		return err
	}

	theme, err := c.getTheme()
	if err != nil {
		return err
	}

	writeOutput := func() {
		c.writeOutput(theme, postmanCollection)
	}

	writeOutput()

	if c.Config.Watch {
		return c.watchThemeFilesChanges(theme, writeOutput)
	}

	return nil
}

func (c *Default) validateUserInput() error {
	if c.Config.CollectionFile == "" {
		return errors.New("You must provide a collection using the -collection flag")
	}
	return nil
}

func (c *Default) getPostmanEnvironment() (environment postman.Environment, err error) {
	if c.Config.EnvironmentFile == "" {
		return
	}

	environment, err = c.EnvironmentBuilder.FromFile(c.Config.EnvironmentFile)
	if err != nil {
		err = fmt.Errorf("Failed to parse environment file: %v", err)
	}

	return
}

func (c *Default) getPostmanCollection(environment postman.Environment) (postman.Collection, error) {
	options := postman.BuilderOptions{
		IgnoredRequestHeaders:  c.Config.IgnoredRequestHeaders.Values,
		IgnoredResponseHeaders: c.Config.IgnoredResponseHeaders.Values,
		EnvironmentVariables:   environment,
	}
	postmanCollection, err := c.CollectionBuilder.FromFile(c.Config.CollectionFile, options)
	if err != nil {
		return postman.Collection{}, fmt.Errorf("Failed to parse collection file: %v", err)
	}

	return postmanCollection, nil
}

func (c *Default) getTheme() (*themes.Theme, error) {
	usedTheme := c.Config.UsedTheme

	theme, err := c.Themes.Open(usedTheme)
	if err == nil {
		return theme, nil
	}

	if err != themes.ErrThemeNotFound {
		return nil, fmt.Errorf("Failed to open the theme: %v", err)
	}

	fmt.Fprintln(c.Config.Out, color.BlueString("Theme '%v' not found, trying to download it...", usedTheme))
	if err := c.Themes.Download(usedTheme); err != nil {
		return nil, err
	}

	return c.Themes.Open(usedTheme)
}

func (c *Default) writeOutput(theme *themes.Theme, collection postman.Collection) {
	outputFile, err := c.createOutputWriter()
	if err != nil {
		fmt.Fprintln(c.Config.Out, color.RedString(err.Error()))
		return
	}
	defer outputFile.Close()

	fmt.Fprint(c.Config.Out, "Generating output... ")
	if err := c.Renderer.Render(outputFile, theme, collection); err != nil {
		fmt.Fprintln(c.Config.Out, color.RedString("FAIL. %v", err))
		return
	}
	fmt.Fprintln(c.Config.Out, color.GreenString("SUCCESS."))
}

func (c *Default) createOutputWriter() (io.WriteCloser, error) {
	if c.Config.OutputFile == "" {
		return nopCloser{c.Config.Out}, nil
	}

	out, err := os.Create(c.Config.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to create output file: %v", err)
	}
	return out, nil
}

func (c *Default) watchThemeFilesChanges(theme *themes.Theme, action func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	go c.executeActionForEachWatcherEvent(watcher, action)

	if err := watcher.Watch(theme.Path); err != nil {
		return fmt.Errorf("Failed to watch theme directory: %v", err)
	}

	c.sleep()

	return nil
}

func (c *Default) executeActionForEachWatcherEvent(watcher *fsnotify.Watcher, action func()) {
	for {
		select {
		case ev := <-watcher.Event:
			if !ev.IsAttrib() {
				action()
			}
		case err := <-watcher.Error:
			fmt.Fprintln(c.Config.Out, color.RedString("FAIL. %v", err))
		}
	}
}

func (c *Default) sleep() {
	select {}
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
