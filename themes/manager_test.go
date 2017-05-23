package themes_test

import (
	"errors"
	"os"
	"path"

	"github.com/aubm/postmanerator/configuration"
	. "github.com/aubm/postmanerator/themes"
	"github.com/aubm/postmanerator/utils"
	. "github.com/aubm/postmanerator/utils/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

const any = mock.Anything

var _ = Describe("Manager", func() {

	var (
		createdTmpThemesDirectory string
		usedConfig                *configuration.Configuration
		mockCloner                *MockGitAgent
		manager                   *Manager
	)

	BeforeEach(func() {
		usedConfig = &configuration.Configuration{}
		mockCloner = &MockGitAgent{}
		manager = &Manager{
			Config: usedConfig,
			Cloner: mockCloner,
		}
	})

	BeforeEach(func() {
		createdTmpThemesDirectory = path.Join(os.TempDir(), uuid.NewV4().String())
		usedConfig.ThemesDirectory = createdTmpThemesDirectory
		must(os.Mkdir(createdTmpThemesDirectory, 0777))
	})

	AfterEach(func() {
		must(os.RemoveAll(createdTmpThemesDirectory))
	})

	Describe("Delete", func() {

		var (
			themeToDelete = "my-theme"
			themePath     string
			returnedError error
		)

		JustBeforeEach(func() {
			returnedError = manager.Delete(themeToDelete)
		})

		BeforeEach(func() {
			themePath = path.Join(usedConfig.ThemesDirectory, themeToDelete)
		})

		Context("when everything is ok", func() {

			BeforeEach(func() {
				must(os.Mkdir(themePath, 0777))
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should delete the theme directory", func() {
				Expect(themePath).NotTo(BeADirectory())
			})

		})

		Context("when the theme directory does not exist", func() {

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

		})

	})

	Describe("List", func() {

		var (
			returnedThemeList []string
			returnedError     error
		)

		JustBeforeEach(func() {
			returnedThemeList, returnedError = manager.List()
		})

		Context("when everything is ok", func() {

			var expectedThemeList = []string{"bar", "default", "foo"}

			BeforeEach(func() {
				for _, themeName := range expectedThemeList {
					themePathToCreate := path.Join(usedConfig.ThemesDirectory, themeName)
					must(os.Mkdir(themePathToCreate, 0777))
				}
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should return the right theme list", func() {
				Expect(returnedThemeList).To(Equal(expectedThemeList))
			})

		})

		Context("when the themes directory does not exist", func() {

			BeforeEach(func() {
				usedConfig.ThemesDirectory = path.Join(usedConfig.ThemesDirectory, "non_existing_dir")
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(ContainSubstring("There was an error while reading the themes directory:"))
			})

			It("should not return a list of themes", func() {
				Expect(returnedThemeList).To(BeZero())
			})

		})

	})

	Describe("Download", func() {

		var (
			themeToDownload    string
			returnedError      error
			cloneErrorToReturn error
		)

		JustBeforeEach(func() {
			mockCloner.On("Clone", any, any).Return(cloneErrorToReturn)
			returnedError = manager.Download(themeToDownload)
		})

		Context("when the theme is a valid http git url", func() {

			BeforeEach(func() {
				themeToDownload = "http://foo.git/bar.git"
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should clone with the right parameters", func() {
				Expect(len(mockCloner.Calls)).To(Equal(1))
				args := mockCloner.Calls[0].Arguments
				Expect(args.Get(0)).To(Equal([]string{"http://foo.git/bar.git"}))
				Expect(args.Get(1)).To(Equal(utils.CloneOptions{TargetDirectory: usedConfig.ThemesDirectory}))
			})

		})

		Context("when the theme is a valid git url", func() {

			BeforeEach(func() {
				themeToDownload = "git@foo.git/bar"
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should clone with the right parameters", func() {
				Expect(len(mockCloner.Calls)).To(Equal(1))
				args := mockCloner.Calls[0].Arguments
				Expect(args.Get(0)).To(Equal([]string{"git@foo.git/bar"}))
				Expect(args.Get(1)).To(Equal(utils.CloneOptions{TargetDirectory: usedConfig.ThemesDirectory}))
			})

		})

		Context("when the clone fails", func() {

			BeforeEach(func() {
				themeToDownload = "git@foo.git/bar"
				cloneErrorToReturn = errors.New("something bad happened!")
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("something bad happened!"))
			})

		})

	})

})
