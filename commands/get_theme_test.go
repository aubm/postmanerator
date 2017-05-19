package commands_test

import (
	"bytes"
	"errors"

	. "github.com/aubm/postmanerator/commands"
	"github.com/aubm/postmanerator/configuration"
	. "github.com/aubm/postmanerator/themes/mocks"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("GetTheme", func() {

	var (
		mockStdOut        *bytes.Buffer
		mockThemesManager *MockThemeManager
		getThemeCommand   *GetTheme
	)

	BeforeEach(func() {
		mockStdOut = new(bytes.Buffer)
		mockThemesManager = &MockThemeManager{}
		mockConfig := &configuration.Configuration{Out: mockStdOut, Args: []string{"themes", "get", "my-custom-theme"}}
		getThemeCommand = &GetTheme{
			Config: mockConfig,
			Themes: mockThemesManager,
		}
	})

	Describe("Is", func() {

		It("should be OK", func() {
			Expect(getThemeCommand.Is("cmd_themes_get")).To(BeTrue())
		})

		It("should be KO", func() {
			Expect(getThemeCommand.Is("cmd_default")).To(BeFalse())
		})

	})

	Describe("Do", func() {

		var returnedError error

		JustBeforeEach(func() {
			returnedError = getThemeCommand.Do()
		})

		Context("when everything is ok", func() {

			BeforeEach(func() {
				mockThemesManager.On("Download", mock.Anything).Return(nil)
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should confirm that the download succeeded in the output", func() {
				Expect(mockStdOut.String()).To(Equal(color.GreenString("Theme successfully downloaded") + "\n"))
			})

			It("should download the right theme", func() {
				Expect(len(mockThemesManager.Calls)).To(Equal(1))
				Expect(mockThemesManager.Calls[0].Arguments.String(0)).To(Equal("my-custom-theme"))
			})

		})

		Context("when no theme is specified", func() {

			BeforeEach(func() {
				mockThemesManager.On("Download", mock.Anything).Return(nil)
				getThemeCommand.Config.Args = getThemeCommand.Config.Args[:2]
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("You must provide the name or the URL of the theme you want to download"))
			})

			It("should not print anything in the output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

		Context("when the download fails", func() {

			var managerError = errors.New("something bad happened")

			BeforeEach(func() {
				mockThemesManager.On("Download", mock.Anything).Return(managerError)
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
