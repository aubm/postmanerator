package commands_test

import (
	"bytes"
	"fmt"
	"os"
	"path"

	. "github.com/aubm/postmanerator/commands"
	"github.com/aubm/postmanerator/configuration"
	"github.com/aubm/postmanerator/postman"
	. "github.com/aubm/postmanerator/postman/mocks"
	"github.com/aubm/postmanerator/themes"
	. "github.com/aubm/postmanerator/themes/mocks"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

const any = mock.Anything

var _ = Describe("Default", func() {

	var (
		mockStdOut             *bytes.Buffer
		outputFilePath         string
		mockThemeManager       *MockThemeManager
		mockThemeRenderer      *MockThemeRenderer
		mockCollectionBuilder  *MockCollectionBuilder
		mockEnvironmentBuilder *MockEnvironmentBuilder
		defaultCommand         *Default
	)

	BeforeEach(func() {
		mockStdOut = new(bytes.Buffer)
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
				collection *postman.Collection
				theme      *themes.Theme
			)

			BeforeEach(func() {
				collection = &postman.Collection{Id: "foo"}
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
						IgnoredRequestHeaders:  postman.HeadersList{"X-Foo", "X-Bar"},
						IgnoredResponseHeaders: postman.HeadersList{"X-Fizz", "X-Buzz"},
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

	})

})
