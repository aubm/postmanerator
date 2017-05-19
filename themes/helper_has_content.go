package themes

import "strings"

func helperHasContent(input string) bool {
	input = strings.Trim(input, " ")
	input = strings.Trim(input, "\n")
	return input != ""
}
