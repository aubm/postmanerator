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

		It("generate a hard coded output", func() {
			usedTheme = &Theme{Files: []string{"tests_data/themes/hard_coded/index.tpl"}}
			expectedOutput = readFileContent("tests_data/themes/hard_coded.out")
		})

		It("generate a simple output", func() {
			usedTheme = &Theme{Files: []string{"tests_data/themes/simple/index.tpl"}}
			expectedOutput = readFileContent("tests_data/themes/simple.out")
		})

	})

})
