package themes

import "gopkg.in/russross/blackfriday.v2"

func helperMarkdown(input string) string {
	return string(blackfriday.Run([]byte(input)))
}
