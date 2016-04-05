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

const THEMES_REPOSITORY string = "https://raw.githubusercontent.com/aubm/postmanerator-themes/master/.gitmodules"

var collectionFile = flag.String("collection", "", "the postman exported collection JSON file")
var usedTheme = flag.String("theme", "default", "the theme to use")
var outputFile = flag.String("output", "", "the output file, default is stdout")
var watch = flag.Bool("watch", false, "automatically regenerate the output when the theme changes")
var localName = flag.String("local-name", "", "the name of the local copy of the downloaded theme")
var ignoredRequestHeaders StringsFlag
var ignoredResponseHeaders StringsFlag

var themesDirectory string
var out *os.File = os.Stdout
var args []string

var emptyErr = errors.New("")

func main() {
	flag.Var(&ignoredResponseHeaders, "ignored-response-headers", "a comma seperated list of ignored response headers")
	flag.Var(&ignoredRequestHeaders, "ignored-request-headers", "a comma seperated list of ignored request headers")
	flag.Parse()

	themesDirectory = os.Getenv("POSTMANERATOR_PATH")
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

	args = flag.Args()

	if len(args) == 0 {
		defaultCommand()
		return
	}

	printHelp := func() {
		documentationURL := "https://github.com/aubm/postmanerator"
		checkAndPrintErr(emptyErr, fmt.Sprintf("Command '%v' not found, please see the documention at %v", strings.Join(args, " "), documentationURL))
	}

	switch args[0] {
	case "themes":
		if len(args) < 2 {
			printHelp()
		}
		switch args[1] {
		case "get":
			getTheme()
		case "delete":
			deleteTheme()
		case "list":
			listThemes()
		default:
			printHelp()
		}
	default:
		printHelp()
	}
}

func defaultCommand() {
	var err error

	if *collectionFile == "" {
		checkAndPrintErr(emptyErr, "You must provide a collection using the -collection flag")
	}

	if *outputFile != "" {
		out, err = os.Create(*outputFile)
		checkAndPrintErr(err, fmt.Sprintf("Failed to create output: %v", err))
		defer out.Close()
	}

	col, err := postman.CollectionFromFile(*collectionFile, postman.CollectionOptions{
		IgnoredRequestHeaders:  postman.HeadersList(ignoredRequestHeaders.values),
		IgnoredResponseHeaders: postman.HeadersList(ignoredResponseHeaders.values),
	})
	checkAndPrintErr(err, fmt.Sprintf("Failed to parse collection file: %v", err))

	col.ExtractStructuresDefinition()

	themePath, err := getThemePath()
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

	if *watch {
		watchDir(themePath, func() { generate(out, themeFiles, col) })
	}
}

func getThemePath() (string, error) {
	themePath, err := theme.GetThemePath(*usedTheme, themesDirectory)
	if err != nil {
		if ok, _ := regexp.MatchString(`\/|\\`, *usedTheme); ok == false {
			fmt.Println(color.BlueString("Theme '%v' not found, trying to download it...", *usedTheme))
			if err := theme.GitClone(*usedTheme, "", THEMES_REPOSITORY, theme.DefaultCloner{themesDirectory}); err == nil {
				return theme.GetThemePath(*usedTheme, themesDirectory)
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

func getTheme() {
	if len(args) < 3 {
		checkAndPrintErr(emptyErr, "You must provide the name or the URL of the theme you want to download")
	}

	err := theme.GitClone(args[2], *localName, THEMES_REPOSITORY, theme.DefaultCloner{themesDirectory})
	checkAndPrintErr(err, "")

	fmt.Println(color.GreenString("Theme successfully downloaded"))
}

func deleteTheme() {
	if len(args) < 3 {
		checkAndPrintErr(emptyErr, "You must provide the name of the theme you want to delete")
	}

	err := theme.Delete(args[2], themesDirectory)
	checkAndPrintErr(err, "")

	fmt.Println(color.GreenString("Theme successfully deleted"))
}

func listThemes() {
	err := theme.ListThemes(os.Stdout, themesDirectory)
	checkAndPrintErr(err, "")
}
