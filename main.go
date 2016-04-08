package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
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
	CollectionFile         string
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
	flag.StringVar(&config.CollectionFile, "collection", "", "the postman exported collection JSON file")
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
		defaultCommand(config)
		return
	}

	switch config.Args[0] {
	case "themes":
		if len(config.Args) >= 2 {
			switch config.Args[1] {
			case "get":
				getTheme(config)
			case "delete":
				deleteTheme(config)
			case "list":
				listThemes(config)
			}
		}
		return
	}

	documentationURL := "https://github.com/aubm/postmanerator"
	checkAndPrintErr(emptyErr, fmt.Sprintf("Command '%v' not found, please see the documention at %v", strings.Join(config.Args, " "), documentationURL))
}

func defaultCommand(config Configuration) {
	var (
		err error
		out *os.File = os.Stdout
	)

	if config.CollectionFile == "" {
		checkAndPrintErr(emptyErr, "You must provide a collection using the -collection flag")
	}

	if config.OutputFile != "" {
		out, err = os.Create(config.OutputFile)
		checkAndPrintErr(err, fmt.Sprintf("Failed to create output: %v", err))
		defer out.Close()
	}

	col, err := postman.CollectionFromFile(config.CollectionFile, postman.CollectionOptions{
		IgnoredRequestHeaders:  postman.HeadersList(config.IgnoredRequestHeaders.values),
		IgnoredResponseHeaders: postman.HeadersList(config.IgnoredResponseHeaders.values),
	})
	checkAndPrintErr(err, fmt.Sprintf("Failed to parse collection file: %v", err))

	col.ExtractStructuresDefinition()

	themePath, err := getThemePath(config)
	if err != nil {
		checkAndPrintErr(emptyErr, "The theme was not found")
	}
	themeFiles, err := theme.ListThemeFiles(themePath)
	if err != nil {
		checkAndPrintErr(err, "")
	}

	defer func() {
		if r := recover(); r != nil {
			color.Red("FAIL. %v\n", r)
		}
	}()
	generate(out, themeFiles, col)

	if config.Watch {
		watchDir(themePath, func() { generate(out, themeFiles, col) })
	}
}

func getThemePath(config Configuration) (string, error) {
	themePath, err := theme.GetThemePath(config.UsedTheme, config.ThemesDirectory)
	if err != nil {
		if ok, _ := regexp.MatchString(`\/|\\`, config.UsedTheme); ok == false {
			fmt.Println(color.BlueString("Theme '%v' not found, trying to download it...", config.UsedTheme))
			if err := theme.GitClone(config.UsedTheme, "", THEMES_REPOSITORY, theme.DefaultCloner{config.ThemesDirectory}); err == nil {
				return theme.GetThemePath(config.UsedTheme, config.ThemesDirectory)
			}
		}
		return "", err
	}
	return themePath, nil
}

func generate(out *os.File, themeFiles []string, col *postman.Collection) {
	fmt.Print("Generating output ... ")
	out.Truncate(0)
	templates := template.Must(template.New("").Funcs(helper.GetFuncMap()).ParseFiles(themeFiles...))
	err := templates.ExecuteTemplate(out, "index.tpl", *col)
	if err != nil {
		color.Red("FAIL. %v\n", err)
		return
	}
	fmt.Print(color.GreenString("SUCCESS.\n"))
}

func watchDir(dir string, action func()) {
	watcher, err := fsnotify.NewWatcher()
	checkAndPrintErr(err, fmt.Sprintf("Failed to create file watcher: %v", err))
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				color.Red("FAIL. %v\n", r)
			}
			watchDir(dir, action)
		}()

		for {
			select {
			case ev := <-watcher.Event:
				if !ev.IsAttrib() {
					action()
				}
			case err := <-watcher.Error:
				log.Printf("FAIL. %v\n", err)
			}
		}
	}()

	err = watcher.Watch(dir)
	checkAndPrintErr(err, fmt.Sprintf("Failed to watch theme directory: %v", err))

	<-done
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

func getTheme(config Configuration) {
	if len(config.Args) < 3 {
		checkAndPrintErr(emptyErr, "You must provide the name or the URL of the theme you want to download")
	}

	err := theme.GitClone(config.Args[2], config.LocalName, THEMES_REPOSITORY, theme.DefaultCloner{config.ThemesDirectory})
	checkAndPrintErr(err, "")

	fmt.Println(color.GreenString("Theme successfully downloaded"))
}

func deleteTheme(config Configuration) {
	if len(config.Args) < 3 {
		checkAndPrintErr(emptyErr, "You must provide the name of the theme you want to delete")
	}

	err := theme.Delete(config.Args[2], config.ThemesDirectory)
	checkAndPrintErr(err, "")

	fmt.Println(color.GreenString("Theme successfully deleted"))
}

func listThemes(config Configuration) {
	err := theme.ListThemes(os.Stdout, config.ThemesDirectory)
	checkAndPrintErr(err, "")
}
