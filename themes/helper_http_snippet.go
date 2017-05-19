package themes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aubm/postmanerator/postman"
)

func helperHttpSnippet(request postman.Request) (httpSnippet string) {
	parsedURL, err := request.ParsedURL()
	if err != nil {
		httpSnippet = err.Error()
		return
	}

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
