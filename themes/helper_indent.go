package themes

import (
	"strings"
)

func helperIndent(spaces int, input string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(input, "\n", "\n"+pad, -1)
}

func helperNIndent(spaces int, input string) string {
	return "\n" + helperIndent(spaces, input)
}
