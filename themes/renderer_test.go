package themes_test

import (
	"bytes"

	"github.com/aubm/postmanerator/postman"
	. "github.com/aubm/postmanerator/themes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Renderer", func() {

	var renderer *Renderer

	BeforeEach(func() {
		renderer = &Renderer{}
	})

	Describe("Render", func() {

		var (
			outputWriter   *bytes.Buffer
			collection     postman.Collection
			expectedOutput string
			usedTheme      *Theme
			returnedError  error
		)

		BeforeEach(func() {
			outputWriter = new(bytes.Buffer)
			collection = exampleCollection
			expectedOutput = ""
			usedTheme = nil
			returnedError = nil
		})

		AfterEach(func() {
			returnedError = renderer.Render(outputWriter, usedTheme, collection)
			Expect(returnedError).To(BeNil())
			Expect(outputWriter.String()).To(Equal(expectedOutput))
		})

		It("should generate a hard coded output", func() {
			usedTheme = &Theme{Files: []string{"tests_data/themes/hard_coded/index.tpl"}}
			expectedOutput = readFileContent("tests_data/themes/hard_coded.out")
		})

		It("should generate a simple output", func() {
			usedTheme = &Theme{Files: []string{"tests_data/themes/simple/index.tpl"}}
			expectedOutput = readFileContent("tests_data/themes/simple.out")
		})

		It("should generate curl snippets", func() {
			usedTheme = &Theme{Files: []string{"tests_data/themes/curl_snippets/index.tpl"}}
			expectedOutput = readFileContent("tests_data/themes/curl_snippets.out")
		})

		It("should generate http snippets", func() {
			usedTheme = &Theme{Files: []string{"tests_data/themes/http_snippets/index.tpl"}}
			expectedOutput = readFileContent("tests_data/themes/http_snippets.out")
		})

	})

})
