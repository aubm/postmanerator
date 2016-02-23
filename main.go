package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"text/template"

	"github.com/aubm/postmanerator/postman"
	"github.com/russross/blackfriday"
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

	templates := template.Must(template.New("").Funcs(template.FuncMap{
		"findRequest":  findRequest,
		"findResponse": findResponse,
		"markdown":     markdown,
		"randomID":     randomID,
		"indentJSON":   indentJSON,
	}).ParseGlob(fmt.Sprintf("./themes/%v/index.tpl", *theme)))
	err = templates.ExecuteTemplate(out, "index.tpl", *col)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findRequest(requests []postman.Request, ID string) *postman.Request {
	for _, r := range requests {
		if r.ID == ID {
			return &r
		}
	}
	return nil
}

func findResponse(req postman.Request, name string) *postman.Response {
	for _, res := range req.Responses {
		if res.Name == name {
			return &res
		}
	}
	return nil
}

func markdown(input string) string {
	return string(blackfriday.MarkdownBasic([]byte(input)))
}

func randomID() int {
	return rand.Intn(999999999)
}

func indentJSON(input string) (string, error) {
	dest := new(bytes.Buffer)
	src := []byte(input)
	err := json.Indent(dest, src, "", "    ")
	return dest.String(), err
}
