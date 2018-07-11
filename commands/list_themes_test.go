package commands_test

import (
	"bytes"
	"errors"

	. "github.com/aubm/postmanerator/commands"
	"github.com/aubm/postmanerator/configuration"
	. "github.com/aubm/postmanerator/themes/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListThemes", func() {

	var (
		mockStdOut        *bytes.Buffer
		mockThemesManager *MockThemeManager
		listThemesCommand *ListThemes
	)

	BeforeEach(func() {
		mockStdOut = new(bytes.Buffer)
		mockThemesManager = &MockThemeManager{}
		mockConfig := &configuration.Configuration{Out: mockStdOut}
		listThemesCommand = &ListThemes{
			Config: mockConfig,
			Themes: mockThemesManager,
		}
	})

	Describe("Is", func() {

		It("should be OK", func() {
			Expect(listThemesCommand.Is("cmd_themes_list")).To(BeTrue())
		})

		It("should be KO", func() {
			Expect(listThemesCommand.Is("cmd_default")).To(BeFalse())
		})

	})

	Describe("Do", func() {

		var returnedError error

		JustBeforeEach(func() {
			returnedError = listThemesCommand.Do()
		})

		Context("when the manager responds with a list of themes", func() {

			BeforeEach(func() {
				mockThemesManager.On("List").Return([]string{"default", "markdown"}, nil)
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should print the themes in the output", func() {
				Expect(mockStdOut.String()).To(Equal("default\nmarkdown\n"))
			})

		})

		Context("when the manager responds with an error", func() {

			var managerError = errors.New("something bad happened")

			BeforeEach(func() {
				mockThemesManager.On("List").Return([]string{}, managerError)
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError).To(Equal(managerError))
			})

			It("should not write anything in the output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

	})

})
