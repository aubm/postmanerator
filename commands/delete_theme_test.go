package commands_test

import (
	"bytes"
	"errors"

	"github.com/fatih/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/srgrn/postmanerator/commands"
	"github.com/srgrn/postmanerator/configuration"
	. "github.com/srgrn/postmanerator/themes/mocks"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("DeleteTheme", func() {

	var (
		mockStdOut         *bytes.Buffer
		mockThemesManager  *MockThemeManager
		deleteThemeCommand *DeleteTheme
	)

	BeforeEach(func() {
		mockStdOut = new(bytes.Buffer)
		mockThemesManager = &MockThemeManager{}
		mockConfig := &configuration.Configuration{Out: mockStdOut, Args: []string{"themes", "delete", "my-custom-theme"}}
		deleteThemeCommand = &DeleteTheme{
			Config: mockConfig,
			Themes: mockThemesManager,
		}
	})

	Describe("Is", func() {

		It("should be OK", func() {
			Expect(deleteThemeCommand.Is("cmd_themes_delete")).To(BeTrue())
		})

		It("should be KO", func() {
			Expect(deleteThemeCommand.Is("cmd_default")).To(BeFalse())
		})

	})

	Describe("Do", func() {

		var returnedError error

		JustBeforeEach(func() {
			returnedError = deleteThemeCommand.Do()
		})

		Context("when everything is ok", func() {

			BeforeEach(func() {
				mockThemesManager.On("Delete", mock.Anything).Return(nil)
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should confirm the deletion in the output", func() {
				Expect(mockStdOut.String()).To(Equal(color.GreenString("Theme successfully deleted") + "\n"))
			})

			It("should delete the right theme", func() {
				Expect(len(mockThemesManager.Calls)).To(Equal(1))
				Expect(mockThemesManager.Calls[0].Arguments.String(0)).To(Equal("my-custom-theme"))
			})

		})

		Context("when no theme is specified", func() {

			BeforeEach(func() {
				mockThemesManager.On("Delete", mock.Anything).Return(nil)
				deleteThemeCommand.Config.Args = deleteThemeCommand.Config.Args[:2]
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("You must provide the name of the theme you want to delete"))
			})

			It("should not print anything in the output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

		Context("when the deletion fails", func() {

			var managerError = errors.New("something bad happened")

			BeforeEach(func() {
				mockThemesManager.On("Delete", mock.Anything).Return(managerError)
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError).To(Equal(managerError))
			})

			It("should not print anything in the output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

	})

})
