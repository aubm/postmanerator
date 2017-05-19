package postman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type CollectionBuilder struct{}

func (b *CollectionBuilder) FromFile(file string, options BuilderOptions) (*Collection, error) {
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

	col.ExtractStructuresDefinition()

	return col, nil
}

type BuilderOptions struct {
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
