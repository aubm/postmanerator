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
func CollectionFromFile(file string) (*Collection, error) {
	col := new(Collection)

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, col)
	if err != nil {
		return nil, err
	}

	return col, nil
}
