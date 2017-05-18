package postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Collection struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Order       []string    `json:"order"`
	Folders     []Folder    `json:"folders"`
	Timestamp   int64       `json:"timestamp"`
	Owner       interface{} `json:"owner"`
	RemoteLink  string      `json:"remoteLink"`
	Public      bool        `json:"public"`
	Requests    []Request   `json:"requests"`
	Structures  []StructureDefinition
}

// CollectionFromFile parses the content of a file and in a new collection
func CollectionFromFile(file string, options CollectionOptions) (*Collection, error) {
	col := new(Collection)

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if options.EnvironmentVariables != nil {
		for k, v := range options.EnvironmentVariables {
			buf = bytes.Replace(buf, []byte(fmt.Sprintf("{{%v}}", k)), []byte(v), -1)
		}
	}

	err = json.Unmarshal(buf, col)
	if err != nil {
		return nil, err
	}

	if len(options.IgnoredRequestHeaders) > 0 {
		for i := 0; i < len(col.Requests); i++ {
			var newRawHeaders []string
			for _, rawHeader := range strings.Split(col.Requests[i].RawHeaders, "\n") {
				if rawHeader == "" {
					continue
				}
				headerName := strings.Split(rawHeader, ":")[0]
				if !options.IgnoredRequestHeaders.Contains(headerName) {
					newRawHeaders = append(newRawHeaders, rawHeader)
				}
			}
			col.Requests[i].RawHeaders = strings.Join(newRawHeaders, "\n")
		}
	}

	if len(options.IgnoredResponseHeaders) > 0 {
		for i := 0; i < len(col.Requests); i++ {
			for j := 0; j < len(col.Requests[i].Responses); j++ {
				newHeaders := []ResponseHeader{}
				for _, header := range col.Requests[i].Responses[j].Headers {
					if options.IgnoredResponseHeaders.Contains(header.Name) == false {
						newHeaders = append(newHeaders, header)
					}
				}
				col.Requests[i].Responses[j].Headers = newHeaders
			}
		}
	}

	return col, nil
}

type CollectionOptions struct {
	IgnoredRequestHeaders  HeadersList
	IgnoredResponseHeaders HeadersList
	EnvironmentVariables   Environment
}

type HeadersList []string

func (list HeadersList) Contains(value string) bool {
	for _, header := range list {
		if header == value {
			return true
		}
	}
	return false
}
