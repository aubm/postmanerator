package themes

import "github.com/russross/blackfriday"

func helperMarkdown(input string) string {
	return string(blackfriday.MarkdownCommon([]byte(input)))
}
