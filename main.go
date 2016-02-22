package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/aubm/postmanerator/postman"
	"github.com/russross/blackfriday"
)

var col *postman.Collection

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

	// Get Postman collection
	col = new(postman.Collection)
	buf, err := ioutil.ReadFile(args[0])
	checkErr(err)

	err = json.Unmarshal(buf, col)
	checkErr(err)

	templates := template.Must(template.New("").Funcs(getFuncMap()).ParseGlob(fmt.Sprintf("./themes/%v/index.tpl", *theme)))
	err = templates.ExecuteTemplate(out, "index.tpl", *col)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"findRequest": func(ID string) (postman.Request, error) {
			for _, r := range col.Requests {
				if r.ID == ID {
					return r, nil
				}
			}
			return postman.Request{}, errors.New("request not found")
		},
		"findResponse": func(req postman.Request, name string) (postman.Response, error) {
			for _, res := range req.Responses {
				if res.Name == name {
					return res, nil
				}
			}
			return postman.Response{}, errors.New("response not found")
		},
		"markdown": func(input string) string {
			return string(blackfriday.MarkdownBasic([]byte(input)))
		},
	}
}
