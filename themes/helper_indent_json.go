package themes

import (
	"bytes"
	"encoding/json"
)

func helperIndentJSON(input string) string {
	dest := new(bytes.Buffer)
	src := []byte(input)
	err := json.Indent(dest, src, "", "    ")
	if err != nil {
		return ""
	}
	return dest.String()
}
