package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/aubm/postmanerator/postman"
	"github.com/aubm/postmanerator/theme/helper"
	"github.com/fatih/color"
	"github.com/howeyc/fsnotify"
)

var theme = flag.String("theme", "markdown_default", "the theme to use")
var outputFile = flag.String("output", "", "the output file, default is stdout")
var watch = flag.Bool("watch", false, "automatically regenerate the output when the theme changes")

var out *os.File = os.Stdout

func main() {
	flag.Parse()

	var err error
	args := flag.Args()

	if len(args) != 1 {
		checkAndPrintErr(errors.New(""), "Missing collection path.")
	}

	if *outputFile != "" {
		out, err = os.Create(*outputFile)
		checkAndPrintErr(err, fmt.Sprintf("Failed to create output: %v", err))
		defer out.Close()
	}

	col, err := postman.CollectionFromFile(args[0])
	checkAndPrintErr(err, fmt.Sprintf("Failed to parse collection file: %v", err))

	col.ExtractStructuresDefinition()

	defer func() {
		if r := recover(); r != nil {
			color.Red("FAIL. %v\n", r)
		}
	}()
	generate(out, col)

	if *watch {
		watchDir(*theme, func() { generate(out, col) })
	}
}

func generate(out *os.File, col *postman.Collection) {
	fmt.Print("Generating output ... ")
	out.Truncate(0)
	templates := template.Must(template.New("").Funcs(helper.GetFuncMap()).ParseGlob(fmt.Sprintf("%v/index.tpl", *theme)))
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
		fmt.Println(color.RedString(msg))
		os.Exit(1)
	}
}
