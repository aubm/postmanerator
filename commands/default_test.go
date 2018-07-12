package commands_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/fatih/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	. "github.com/srgrn/postmanerator/commands"
	"github.com/srgrn/postmanerator/configuration"
	"github.com/srgrn/postmanerator/postman"
	. "github.com/srgrn/postmanerator/postman/mocks"
	"github.com/srgrn/postmanerator/themes"
	. "github.com/srgrn/postmanerator/themes/mocks"
	"github.com/stretchr/testify/mock"
)

const any = mock.Anything

var _ = Describe("Default", func() {

	var (
		mockStdOut             *bytes.Buffer
		someBadError           error
		outputFilePath         string
		mockThemeManager       *MockThemeManager
		mockThemeRenderer      *MockThemeRenderer
		mockCollectionBuilder  *MockCollectionBuilder
		mockEnvironmentBuilder *MockEnvironmentBuilder
		defaultCommand         *Default
	)

	BeforeEach(func() {
		mockStdOut = new(bytes.Buffer)
		someBadError = errors.New("something bad happened!")
		outputFilePath = path.Join(os.TempDir(), fmt.Sprintf("postmanerator-generated-test-output-%s.out", uuid.NewV4().String()))
		mockThemeManager = &MockThemeManager{}
		mockThemeRenderer = &MockThemeRenderer{}
		mockCollectionBuilder = &MockCollectionBuilder{}
		mockEnvironmentBuilder = &MockEnvironmentBuilder{}
		mockConfig := &configuration.Configuration{
			Out:            mockStdOut,
			UsedTheme:      "default",
			CollectionFile: "awesome-collection.json",
			OutputFile:     outputFilePath,
		}
		defaultCommand = &Default{
			Config:             mockConfig,
			Themes:             mockThemeManager,
			Renderer:           mockThemeRenderer,
			CollectionBuilder:  mockCollectionBuilder,
			EnvironmentBuilder: mockEnvironmentBuilder,
		}
	})

	AfterEach(func() {
		os.Remove(outputFilePath)
	})

	Describe("Is", func() {

		It("should be OK", func() {
			Expect(defaultCommand.Is("cmd_default")).To(BeTrue())
		})

		It("should be KO", func() {
			Expect(defaultCommand.Is("cmd_themes_delete")).To(BeFalse())
		})

	})

	Describe("Do", func() {

		var returnedError error

		JustBeforeEach(func() {
			returnedError = defaultCommand.Do()
		})

		Context("when everything is ok", func() {

			var (
				collection postman.Collection
				theme      *themes.Theme
			)

			BeforeEach(func() {
				collection = postman.Collection{Name: "foo"}
				theme = &themes.Theme{Name: "foo"}
				mockCollectionBuilder.On("FromFile", any, any).Return(collection, nil)
				mockThemeManager.On("Open", any).Return(theme, nil)
				mockThemeRenderer.On("Render", any, any, any).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprint(args.Get(0).(*os.File), "some contents")
				})
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should not try to parse an environment file", func() {
				Expect(len(mockEnvironmentBuilder.Calls)).To(Equal(0))
			})

			It("should build the right collection file with the appropriate options", func() {
				Expect(len(mockCollectionBuilder.Calls)).To(Equal(1))
				args := mockCollectionBuilder.Calls[0].Arguments
				Expect(args.String(0)).To(Equal("awesome-collection.json"))
				Expect(args.Get(1)).To(Equal(postman.BuilderOptions{}))
			})

			It("should open the right theme", func() {
				Expect(len(mockThemeManager.Calls)).To(Equal(1))
				args := mockThemeManager.Calls[0].Arguments
				Expect(args.String(0)).To(Equal("default"))
			})

			It("should render the right theme and the right collection", func() {
				Expect(len(mockThemeRenderer.Calls)).To(Equal(1))
				args := mockThemeRenderer.Calls[0].Arguments
				Expect(args.Get(1)).To(Equal(theme))
				Expect(args.Get(2)).To(Equal(collection))
			})

			It("should produce contents in the right output file", func() {
				Expect(readFileContents(outputFilePath)).To(Equal("some contents"))
			})

			It("should produce the right command output", func() {
				Expect(mockStdOut.String()).To(Equal("Generating output... " + color.GreenString("SUCCESS.") + "\n"))
			})

			Context("and custom config parameters for the collection parsing", func() {

				BeforeEach(func() {
					defaultCommand.Config.IgnoredRequestHeaders = configuration.StringsFlag{Values: []string{"X-Foo", "X-Bar"}}
					defaultCommand.Config.IgnoredResponseHeaders = configuration.StringsFlag{Values: []string{"X-Fizz", "X-Buzz"}}
				})

				It("should propagate the options to the collection builder", func() {
					args := mockCollectionBuilder.Calls[0].Arguments
					Expect(args.Get(1)).To(Equal(postman.BuilderOptions{
						IgnoredRequestHeaders:  []string{"X-Foo", "X-Bar"},
						IgnoredResponseHeaders: []string{"X-Fizz", "X-Buzz"},
					}))
				})

			})

			Context("and using a custom environment", func() {

				var environment postman.Environment

				BeforeEach(func() {
					defaultCommand.Config.EnvironmentFile = "awesome-environment.json"
					environment = postman.Environment{"foo": "bar"}
					mockEnvironmentBuilder.On("FromFile", any).Return(environment, nil)
				})

				It("should build the right environment file", func() {
					Expect(len(mockEnvironmentBuilder.Calls)).To(Equal(1))
					Expect(mockEnvironmentBuilder.Calls[0].Arguments.String(0)).To(Equal("awesome-environment.json"))
				})

				It("should propagate the environment to the collection builder", func() {
					args := mockCollectionBuilder.Calls[0].Arguments
					Expect(args.Get(1)).To(Equal(postman.BuilderOptions{
						EnvironmentVariables: environment,
					}))
				})

			})

			Context("and the output file already exists and has contents", func() {

				BeforeEach(func() {
					putFileContents(outputFilePath, "existing contents")
				})

				It("should truncate the existing contents", func() {
					Expect(readFileContents(outputFilePath)).To(Equal("some contents"))
				})

			})

		})

		Context("when no collection is provided", func() {

			BeforeEach(func() {
				defaultCommand.Config.CollectionFile = ""
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("You must provide a collection using the -collection flag"))
			})

			It("should not produce any command output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

		Context("when parsing the environment file fails", func() {

			BeforeEach(func() {
				defaultCommand.Config.EnvironmentFile = "invalid-environment.json"
				mockEnvironmentBuilder.On("FromFile", any).Return(postman.Environment{}, someBadError)
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("Failed to parse environment file: something bad happened!"))
			})

			It("should not produce any command output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

		Context("when parsing the collection file fails", func() {

			BeforeEach(func() {
				mockCollectionBuilder.On("FromFile", any, any).Return(postman.Collection{}, someBadError)
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("Failed to parse collection file: something bad happened!"))
			})

			It("should not produce any command output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

		Context("when opening the theme fails", func() {

			var collection postman.Collection

			BeforeEach(func() {
				collection = postman.Collection{Name: "foo"}
				mockCollectionBuilder.On("FromFile", any, any).Return(collection, nil)
				mockThemeManager.On("Open", any).Return(&themes.Theme{}, someBadError)
			})

			It("should return an error", func() {
				Expect(returnedError).ToNot(BeNil())
				Expect(returnedError.Error()).To(Equal("Failed to open the theme: something bad happened!"))
			})

			It("should not produce any command output", func() {
				Expect(mockStdOut.String()).To(BeZero())
			})

		})

		Context("when the theme is not found locally", func() {

			var (
				collection postman.Collection
				theme      *themes.Theme
			)

			BeforeEach(func() {
				defaultCommand.Config.UsedTheme = "custom_theme"
				collection = postman.Collection{Name: "foo"}
				theme = &themes.Theme{Name: "custom_theme"}
				mockCollectionBuilder.On("FromFile", any, any).Return(collection, nil)
				mockThemeManager.On("Open", any).Return(&themes.Theme{}, themes.ErrThemeNotFound).Once()
				mockThemeManager.On("Download", any).Return(nil).Once()
				mockThemeManager.On("Open", any).Return(theme, nil).Once()
				mockThemeRenderer.On("Render", any, any, any).Return(nil)
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should try to open, download and try to open again the right theme", func() {
				Expect(len(mockThemeManager.Calls)).To(Equal(3))
				for _, c := range mockThemeManager.Calls {
					Expect(c.Arguments.String(0)).To(Equal("custom_theme"))
				}
			})

			It("should produce the right command output", func() {
				expectedOutput := color.BlueString("Theme 'custom_theme' not found, trying to download it...") + "\n" +
					"Generating output... " + color.GreenString("SUCCESS.") + "\n"
				Expect(mockStdOut.String()).To(Equal(expectedOutput))
			})

		})

		Context("when the theme does not exist", func() {

			var collection postman.Collection

			BeforeEach(func() {
				defaultCommand.Config.UsedTheme = "custom_theme"
				collection = postman.Collection{Name: "foo"}
				mockCollectionBuilder.On("FromFile", any, any).Return(collection, nil)
				mockThemeManager.On("Open", any).Return(&themes.Theme{}, themes.ErrThemeNotFound).Once()
				mockThemeManager.On("Download", any).Return(someBadError).Once()
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("something bad happened!"))
			})

			It("should try to open, download and try to open again the right theme", func() {
				Expect(len(mockThemeManager.Calls)).To(Equal(2))
				for _, c := range mockThemeManager.Calls {
					Expect(c.Arguments.String(0)).To(Equal("custom_theme"))
				}
			})

			It("should produce the right command output", func() {
				Expect(mockStdOut.String()).To(Equal(color.BlueString("Theme 'custom_theme' not found, trying to download it...") + "\n"))
			})

		})

		Context("when no output file is specified", func() {

			var (
				collection postman.Collection
				theme      *themes.Theme
			)

			BeforeEach(func() {
				defaultCommand.Config.OutputFile = ""
				collection = postman.Collection{Name: "foo"}
				theme = &themes.Theme{Name: "foo"}
				mockCollectionBuilder.On("FromFile", any, any).Return(collection, nil)
				mockThemeManager.On("Open", any).Return(theme, nil).Once()
				mockThemeRenderer.On("Render", any, any, any).Return(nil).Run(func(args mock.Arguments) {
					fmt.Fprint(args.Get(0).(io.Writer), "some contents")
				})
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should write in the standard output", func() {
				Expect(mockStdOut.String()).To(Equal("Generating output... some contents" + color.GreenString("SUCCESS.") + "\n"))
			})

		})

		Context("when the output file can not be created", func() {

			var (
				collection postman.Collection
				theme      *themes.Theme
			)

			BeforeEach(func() {
				collection = postman.Collection{Name: "foo"}
				theme = &themes.Theme{Name: "foo"}
				outputFilePath = path.Join(os.TempDir(), fmt.Sprintf("postmanerator-%s/embedded/dir", uuid.NewV4().String()))
				defaultCommand.Config.OutputFile = outputFilePath
				mockCollectionBuilder.On("FromFile", any, any).Return(collection, nil)
				mockThemeManager.On("Open", any).Return(theme, nil)
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should produce the right command output", func() {
				Expect(mockStdOut.String()).To(ContainSubstring("Failed to create output file: open"))
			})

		})

		Context("when rendering fails", func() {

			var (
				collection postman.Collection
				theme      *themes.Theme
			)

			BeforeEach(func() {
				collection = postman.Collection{Name: "foo"}
				theme = &themes.Theme{Name: "foo"}
				mockCollectionBuilder.On("FromFile", any, any).Return(collection, nil)
				mockThemeManager.On("Open", any).Return(theme, nil)
				mockThemeRenderer.On("Render", any, any, any).Return(someBadError)
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should produce the right command output", func() {
				Expect(mockStdOut.String()).To(Equal("Generating output... " + color.RedString("FAIL. something bad happened!") + "\n"))
			})

		})

	})

})
