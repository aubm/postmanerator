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
)

var theme = flag.String("theme", "markdown_default", "the theme to use")
var outputFile = flag.String("output", "", "the output file, default is stdout")

var out *os.File = os.Stdout

func main() {
	flag.Parse()

	var err error
	args := flag.Args()

	if len(args) != 1 {
		checkErr(errors.New("Missing collection path"))
	}

	if *outputFile != "" {
		out, err = os.Create(*outputFile)
		checkErr(err)
		defer out.Close()
	}

	col, err := postman.CollectionFromFile(args[0])
	checkErr(err)

	col.ExtractStructuresDefinition()

	templates := template.Must(template.New("").Funcs(helper.GetFuncMap()).ParseGlob(fmt.Sprintf("%v/index.tpl", *theme)))
	err = templates.ExecuteTemplate(out, "index.tpl", *col)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
