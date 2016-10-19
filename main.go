package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"regexp"
	"strings"
	"text/template"

	"github.com/aubm/postmanerator/postman"
	"github.com/aubm/postmanerator/theme"
	"github.com/aubm/postmanerator/theme/helper"
	"github.com/fatih/color"
	"github.com/howeyc/fsnotify"
)

type Configuration struct {
	Out                    io.Writer
	CollectionFile         string
	EnvironmentFile        string
	UsedTheme              string
	OutputFile             string
	Watch                  bool
	LocalName              string
	IgnoredRequestHeaders  StringsFlag
	IgnoredResponseHeaders StringsFlag
	ThemesDirectory        string
	Args                   []string
}

const THEMES_REPOSITORY string = "https://raw.githubusercontent.com/aubm/postmanerator-themes/master/.gitmodules"

var (
	emptyErr = errors.New("")
	config   Configuration
)

func init() {
	config.Out = os.Stdout
	flag.StringVar(&config.CollectionFile, "collection", "", "the postman exported collection JSON file")
	flag.StringVar(&config.EnvironmentFile, "environment", "", "the postman exported environment JSON file")
	flag.StringVar(&config.UsedTheme, "theme", "default", "the theme to use")
	flag.StringVar(&config.OutputFile, "output", "", "the output file, default is stdout")
	flag.BoolVar(&config.Watch, "watch", false, "automatically regenerate the output when the theme changes")
	flag.StringVar(&config.LocalName, "local-name", "", "the name of the local copy of the downloaded theme")
	flag.Var(&config.IgnoredResponseHeaders, "ignored-response-headers", "a comma seperated list of ignored response headers")
	flag.Var(&config.IgnoredRequestHeaders, "ignored-request-headers", "a comma seperated list of ignored request headers")
	flag.Parse()

	themesDirectory := os.Getenv("POSTMANERATOR_PATH")
	if themesDirectory == "" {
		var usrHomeDir string
		usr, err := user.Current()
		if err != nil {
			if usrHomeDir = os.Getenv("HOME"); usrHomeDir == "" {
				if usrHomeDir = os.Getenv("USERPROFILE"); usrHomeDir == "" {
					checkAndPrintErr(err, `An error occured while trying to determine which directory to use for themes.
As a workaround, you can define the POSTMANERATOR_PATH environement variable.
Please consult the documentation here https://github.com/aubm/postmanerator and feel free to submit an issue.`)
				}
			}
		} else {
			usrHomeDir = usr.HomeDir
		}
		themesDirectory = fmt.Sprintf("%v/.postmanerator", usrHomeDir)
	}
	themesDirectory += "/themes"
	if _, err := os.Stat(themesDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(themesDirectory, 0777)
		checkAndPrintErr(err, fmt.Sprintf("Failed to create themes directory: %v", err))
	}
	config.ThemesDirectory = themesDirectory

	config.Args = flag.Args()
}

func main() {
	if len(config.Args) == 0 {
		if err := defaultCommand(config); err != nil {
			checkAndPrintErr(err, "")
		}
		return
	}

	switch config.Args[0] {
	case "themes":
		var err error
		if len(config.Args) >= 2 {
			switch config.Args[1] {
			case "get":
				err = getTheme(config)
			case "delete":
				err = deleteTheme(config)
			case "list":
				err = listThemes(config)
			}
		}
		checkAndPrintErr(err, "")
		return
	}

	documentationURL := "https://github.com/aubm/postmanerator"
	checkAndPrintErr(emptyErr, fmt.Sprintf("Command '%v' not found, please see the documention at %v", strings.Join(config.Args, " "), documentationURL))
}

func defaultCommand(config Configuration) error {
	var (
		err error
		out *os.File = os.Stdout
	)

	if config.CollectionFile == "" {
		return errors.New("You must provide a collection using the -collection flag")
	}

	var env map[string]string
	if config.EnvironmentFile != "" {
		env, err = postman.EnvironmentFromFile(config.EnvironmentFile)
		if err != nil {
			return fmt.Errorf("Failed to open environment file: %v", err)
		}
	}

	if config.OutputFile != "" {
		out, err = os.Create(config.OutputFile)
		if err != nil {
			return fmt.Errorf("Failed to create output: %v", err)
		}
		defer out.Close()
	}

	col, err := postman.CollectionFromFile(config.CollectionFile, postman.CollectionOptions{
		IgnoredRequestHeaders:  postman.HeadersList(config.IgnoredRequestHeaders.values),
		IgnoredResponseHeaders: postman.HeadersList(config.IgnoredResponseHeaders.values),
		EnvironmentVariables:   env,
	})
	if err != nil {
		return fmt.Errorf("Failed to parse collection file: %v", err)
	}

	col.ExtractStructuresDefinition()

	themePath, err := getThemePath(config)
	if err != nil {
		return errors.New("The theme was not found")
	}
	themeFiles, err := theme.ListThemeFiles(themePath)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			config.Out.Write([]byte(color.RedString("FAIL. %v\n", r)))
		}
	}()
	generate(config.Out, out, themeFiles, col)

	if config.Watch {
		return watchDir(config.Out, themePath, func() { generate(config.Out, out, themeFiles, col) })
	}

	return nil
}

func getThemePath(config Configuration) (string, error) {
	themePath, err := theme.GetThemePath(config.UsedTheme, config.ThemesDirectory)
	if err != nil {
		if ok, _ := regexp.MatchString(`\/|\\`, config.UsedTheme); ok == false {
			config.Out.Write([]byte(color.BlueString("Theme '%v' not found, trying to download it...\n", config.UsedTheme)))
			if err := theme.GitClone(config.UsedTheme, "", THEMES_REPOSITORY, theme.DefaultCloner{config.ThemesDirectory}); err == nil {
				return theme.GetThemePath(config.UsedTheme, config.ThemesDirectory)
			}
		}
		return "", err
	}
	return themePath, nil
}

func generate(out io.Writer, f *os.File, themeFiles []string, col *postman.Collection) {
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

func watchDir(out io.Writer, dir string, action func()) error {
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
			watchDir(out, dir, action)
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

func checkAndPrintErr(err error, msg string) {
	if err != nil {
		if msg == "" {
			msg = err.Error()
		}
		fmt.Println(color.RedString(msg))
		os.Exit(1)
	}
}

type StringsFlag struct {
	values []string
}

func (sf StringsFlag) String() string {
	return fmt.Sprint(sf.values)
}

func (sf *StringsFlag) Set(value string) error {
	if value != "" {
		sf.values = strings.Split(value, ",")
	}
	return nil
}

func getTheme(config Configuration) error {
	if len(config.Args) < 3 {
		return errors.New("You must provide the name or the URL of the theme you want to download")
	}

	if err := theme.GitClone(config.Args[2], config.LocalName, THEMES_REPOSITORY, theme.DefaultCloner{config.ThemesDirectory}); err != nil {
		return err
	}

	config.Out.Write([]byte(color.GreenString("Theme successfully downloaded")))
	return nil
}

func deleteTheme(config Configuration) error {
	if len(config.Args) < 3 {
		return errors.New("You must provide the name of the theme you want to delete")
	}

	if err := theme.Delete(config.Args[2], config.ThemesDirectory); err != nil {
		return err
	}

	config.Out.Write([]byte(color.GreenString("Theme successfully deleted")))
	return nil
}

func listThemes(config Configuration) error {
	return theme.ListThemes(os.Stdout, config.ThemesDirectory)
}
