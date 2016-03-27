package postman

import (
	"encoding/json"
	"io/ioutil"
)

type Collection struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Order       []string    `json:"order"`
	Folders     []Folder    `json:"folders"`
	Timestamp   int         `json:"timestamp"`
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

	err = json.Unmarshal(buf, col)
	if err != nil {
		return nil, err
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
	IgnoredResponseHeaders HeadersList
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
