package configuration

import (
	"fmt"
	"strings"
)

type StringsFlag struct {
	Values []string
}

func (sf StringsFlag) String() string {
	return fmt.Sprint(sf.Values)
}

func (sf *StringsFlag) Set(value string) error {
	if value != "" {
		sf.Values = strings.Split(value, ",")
	}
	return nil
}
