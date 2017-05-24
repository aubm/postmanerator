package themes

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/aubm/postmanerator/postman"
)

func helperHttpSnippet(request postman.Request) (httpSnippet string) {
	parsedURL, err := parsedURL(request.URL)
	if err != nil {
		httpSnippet = err.Error()
		return
	}

	httpSnippet += fmt.Sprintf(`%v %v HTTP/1.1
Host: %v`, request.Method, parsedURL.RequestURI(), parsedURL.Host)

	for _, header := range request.Headers {
		httpSnippet += fmt.Sprintf("\n%v: %v", header.Name, header.Value)
	}

	if ok, _ := regexp.MatchString("POST|PUT|PATCH|DELETE", request.Method); ok == false {
		return
	}

	if request.PayloadType == "raw" && request.PayloadRaw != "" {
		httpSnippet += fmt.Sprintf("\n\n%v", request.PayloadRaw)
		return
	}

	if len(request.PayloadParams) <= 0 {
		return
	}

	if request.PayloadType == "urlencoded" {
		var dataList []string
		for _, data := range request.PayloadParams {
			dataList = append(dataList, fmt.Sprintf("%v=%v", data.Key, data.Value))
		}
		httpSnippet += fmt.Sprintf(`
Content-Type: application/x-www-form-urlencoded

%v`, strings.Join(dataList, "&"))
	}

	if request.PayloadType == "params" {
		boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
		httpSnippet += fmt.Sprintf(`
Content-Type: multipart/form-data; boundary=%v

%v`, boundary, boundary)
		for _, data := range request.PayloadParams {
			httpSnippet += fmt.Sprintf(`
Content-Disposition: form-data; name="%v"

%v
%v`, data.Key, data.Value, boundary)
		}
	}

	return
}

func parsedURL(rawUrl string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url\n%v", err)
	}
	return parsedURL, nil
}
