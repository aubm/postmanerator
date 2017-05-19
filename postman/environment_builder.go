package postman

import (
	"encoding/json"
	"os"
)

type EnvironmentBuilder struct{}

func (b *EnvironmentBuilder) FromFile(file string) (Environment, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	envExport := new(environmentExport)

	err = json.NewDecoder(f).Decode(envExport)
	if err != nil {
		return nil, err
	}

	env := map[string]string{}
	for _, v := range envExport.Values {
		env[v.Key] = v.Value
	}

	return env, nil
}

type environmentExport struct {
	Values []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"values"`
}
