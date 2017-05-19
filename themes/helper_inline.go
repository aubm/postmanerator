package themes

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func helperInline(file string) (string, error) {
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
