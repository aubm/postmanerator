package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"text/template"

	"github.com/aubm/postmanerator/postman"
	"github.com/russross/blackfriday"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"findRequest":  findRequest,
		"findResponse": findResponse,
		"markdown":     markdown,
		"randomID":     randomID,
		"indentJSON":   indentJSON,
		"curlSnippet":  curlSnippet,
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

func curlSnippet(request postman.Request) string {
	var formattedHeaders string
	for _, header := range request.Headers() {
		formattedHeaders += fmt.Sprintf(`-H "%v: %v" `, header.Name, header.Value)
	}
	if formattedHeaders == "" {
		formattedHeaders = " "
	}
	return fmt.Sprintf(`curl -X %v %v-d '%v' "%v"`, request.Method, formattedHeaders, request.RawModeData, request.URL)
}
