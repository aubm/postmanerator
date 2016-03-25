package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/aubm/postmanerator/postman"
	"github.com/russross/blackfriday"
)

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"findRequest":  findRequest,
		"findResponse": findResponse,
		"markdown":     markdown,
		"indentJSON":   indentJSON,
		"curlSnippet":  curlSnippet,
		"httpSnippet":  httpSnippet,
		"inline":       inline,
		"slugify":      slugify,
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

func indentJSON(input string) (string, error) {
	dest := new(bytes.Buffer)
	src := []byte(input)
	err := json.Indent(dest, src, "", "    ")
	return dest.String(), err
}

func curlSnippet(request postman.Request) string {
	var curlSnippet string
	payloadReady, _ := regexp.Compile("POST|PUT|PATCH|DELETE")
	curlSnippet += fmt.Sprintf("curl -X %v", request.Method)

	if payloadReady.MatchString(request.Method) {
		if request.DataMode == "urlencoded" {
			curlSnippet += ` -H "Content-Type: application/x-www-form-urlencoded"`
		} else if request.DataMode == "params" {
			curlSnippet += ` -H "Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW"`
		}
	}

	for _, header := range request.Headers() {
		curlSnippet += fmt.Sprintf(` -H "%v: %v"`, header.Name, header.Value)
	}

	if payloadReady.MatchString(request.Method) {
		if request.DataMode == "raw" && request.RawModeData != "" {
			curlSnippet += fmt.Sprintf(` -d '%v'`, request.RawModeData)
		} else if len(request.Data) > 0 {
			if request.DataMode == "urlencoded" {
				var dataList []string
				for _, data := range request.Data {
					dataList = append(dataList, fmt.Sprintf("%v=%v", data.Key, data.Value))
				}
				curlSnippet += fmt.Sprintf(` -d "%v"`, strings.Join(dataList, "&"))
			} else if request.DataMode == "params" {
				for _, data := range request.Data {
					curlSnippet += fmt.Sprintf(` -F "%v=%v"`, data.Key, data.Value)
				}
			}
		}
	}

	curlSnippet += fmt.Sprintf(` "%v"`, request.URL)
	return curlSnippet
}

func httpSnippet(request postman.Request) (httpSnippet string) {
	parsedURL := request.ParsedURL()
	httpSnippet += fmt.Sprintf(`%v %v HTTP/1.1
Host: %v`, request.Method, parsedURL.RequestURI(), parsedURL.Host)

	for _, header := range request.Headers() {
		httpSnippet += fmt.Sprintf("\n%v: %v", header.Name, header.Value)
	}

	if ok, _ := regexp.MatchString("POST|PUT|PATCH|DELETE", request.Method); ok == false {
		return
	}

	if request.DataMode == "raw" && request.RawModeData != "" {
		httpSnippet += fmt.Sprintf("\n\n%v", request.RawModeData)
		return
	}

	if len(request.Data) <= 0 {
		return
	}

	if request.DataMode == "urlencoded" {
		var dataList []string
		for _, data := range request.Data {
			dataList = append(dataList, fmt.Sprintf("%v=%v", data.Key, data.Value))
		}
		httpSnippet += fmt.Sprintf(`
Content-Type: application/x-www-form-urlencoded

%v`, strings.Join(dataList, "&"))
	}

	if request.DataMode == "params" {
		boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
		httpSnippet += fmt.Sprintf(`
Content-Type: multipart/form-data; boundary=%v

%v`, boundary, boundary)
		for _, data := range request.Data {
			httpSnippet += fmt.Sprintf(`
Content-Disposition: form-data; name="%v"

%v
%v`, data.Key, data.Value, boundary)
		}
	}

	return
}

func inline(file string) (string, error) {
	resp, err := http.Get(file)
	if err != nil {
		return "", fmt.Errorf("Failed to fetch URL %v: %v", file, err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read HTTP response for URL %v: %v", file, err)
	}
	return string(content), nil
}

func slugify(label string) string {
	re := regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(label), "-"), "-")
}
