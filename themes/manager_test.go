package themes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/srgrn/postmanerator/configuration"
	. "github.com/srgrn/postmanerator/themes"
	"github.com/srgrn/postmanerator/utils"
	. "github.com/srgrn/postmanerator/utils/mocks"
	"github.com/stretchr/testify/mock"
)

const any = mock.Anything

var _ = Describe("Manager", func() {

	var (
		themesRepositoryServer            *httptest.Server
		themesRepositoryGeneratedRequests []*http.Request
		nbFailedRequests                  int
		nbRequests                        int
		createdTmpThemesDirectory         string
		usedConfig                        *configuration.Configuration
		mockCloner                        *MockGitAgent
		manager                           *Manager
	)

	BeforeEach(func() {
		nbFailedRequests = 0
		nbRequests = 0
		themesRepositoryGeneratedRequests = nil
		themesRepositoryServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nbRequests++
			themesRepositoryGeneratedRequests = append(themesRepositoryGeneratedRequests, r)
			if nbRequests > nbFailedRequests {
				http.ServeFile(w, r, "tests_data/themes_list.txt")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))
	})

	BeforeEach(func() {
		usedConfig = &configuration.Configuration{ThemesRepository: themesRepositoryServer.URL}
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

		BeforeEach(func() {
			themeToDownload = ""
			cloneErrorToReturn = nil
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

		Context("when the user provides a custom local name for the name", func() {

			BeforeEach(func() {
				themeToDownload = "git@foo.git/bar"
				manager.Config.ThemeLocalName = "my-custom-name"
			})

			It("should not return an error", func() {
				Expect(returnedError).To(BeNil())
			})

			It("should clone with the right parameters", func() {
				Expect(len(mockCloner.Calls)).To(Equal(1))
				args := mockCloner.Calls[0].Arguments
				Expect(args.Get(0)).To(Equal([]string{"git@foo.git/bar", "my-custom-name"}))
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

		Context("when the theme is not a git url", func() {

			BeforeEach(func() {
				themeToDownload = "default"
			})

			Context("and everything is ok", func() {

				It("should not return an error", func() {
					Expect(returnedError).To(BeNil())
				})

				It("should generate one request to themes repository", func() {
					Expect(len(themesRepositoryGeneratedRequests)).To(Equal(1))
					Expect(themesRepositoryGeneratedRequests[0].Method).To(Equal(http.MethodGet))
				})

				It("should clone the repository with the default local name", func() {
					Expect(len(mockCloner.Calls)).To(Equal(1))
					args := mockCloner.Calls[0].Arguments
					Expect(args.Get(0)).To(Equal([]string{"https://github.com/aubm/postmanerator-default-theme.git", "default"}))
					Expect(args.Get(1)).To(Equal(utils.CloneOptions{TargetDirectory: usedConfig.ThemesDirectory}))
				})

			})

			Context("and the themes repository request fails two times and finally succeeds on third attempt", func() {

				BeforeEach(func() {
					nbFailedRequests = 2
				})

				It("should not return an error", func() {
					Expect(returnedError).To(BeNil())
				})

				It("should generate three request to themes repository", func() {
					Expect(len(themesRepositoryGeneratedRequests)).To(Equal(3))
				})

				It("should clone the repository with the default local name", func() {
					Expect(len(mockCloner.Calls)).To(Equal(1))
					args := mockCloner.Calls[0].Arguments
					Expect(args.Get(0)).To(Equal([]string{"https://github.com/aubm/postmanerator-default-theme.git", "default"}))
					Expect(args.Get(1)).To(Equal(utils.CloneOptions{TargetDirectory: usedConfig.ThemesDirectory}))
				})

			})

			Context("and the themes repository keeps responding with 500", func() {

				BeforeEach(func() {
					nbFailedRequests = 3
				})

				It("should return an error", func() {
					Expect(returnedError).NotTo(BeNil())
					Expect(returnedError.Error()).To(Equal("failed to download the theme list: themes repository responded with status 500"))
				})

				It("should generate three request to themes repository", func() {
					Expect(len(themesRepositoryGeneratedRequests)).To(Equal(3))
				})

				It("should not clone anything", func() {
					Expect(len(mockCloner.Calls)).To(Equal(0))
				})

			})

			Context("when the themes repository is down", func() {

				BeforeEach(func() {
					themesRepositoryServer.Close()
				})

				It("should return an error", func() {
					Expect(returnedError).NotTo(BeNil())
					Expect(returnedError.Error()).To(ContainSubstring("failed to download the theme list:"))
				})

				It("should not clone anything", func() {
					Expect(len(mockCloner.Calls)).To(Equal(0))
				})

			})

			Context("when the theme is not referenced", func() {

				BeforeEach(func() {
					themeToDownload = "some-unreferenced-theme"
				})

				It("should return an error", func() {
					Expect(returnedError).NotTo(BeNil())
					Expect(returnedError.Error()).To(Equal("Theme not found"))
				})

				It("should generate one request to themes repository", func() {
					Expect(len(themesRepositoryGeneratedRequests)).To(Equal(1))
				})

				It("should not clone anything", func() {
					Expect(len(mockCloner.Calls)).To(Equal(0))
				})

			})

		})

	})

	Describe("Open", func() {

		var (
			themeToOpen   string
			returnedTheme *Theme
			returnedError error
		)

		JustBeforeEach(func() {
			returnedTheme, returnedError = manager.Open(themeToOpen)
		})

		BeforeEach(func() {
			themeToOpen = ""
			returnedTheme = nil
			returnedError = nil
		})

		BeforeEach(func() {
			must(os.Mkdir(path.Join(createdTmpThemesDirectory, "default"), 0777))
			must(os.Create(path.Join(createdTmpThemesDirectory, "default", "index.tpl")))
			must(os.Create(path.Join(createdTmpThemesDirectory, "default", "menu.tpl")))
			must(os.Create(path.Join(createdTmpThemesDirectory, "default", "theme.css")))
			must(os.Create(path.Join(createdTmpThemesDirectory, "invalid-theme")))
		})

		Context("when the theme exist", func() {

			AfterEach(func() {
				Expect(returnedError).To(BeNil())
			})

			Context("given just the theme name", func() {

				BeforeEach(func() {
					themeToOpen = "default"
				})

				It("should return the right theme", func() {
					Expect(returnedTheme).To(Equal(&Theme{
						Name: "default",
						Path: path.Join(createdTmpThemesDirectory, "default"),
						Files: []string{
							path.Join(createdTmpThemesDirectory, "default", "index.tpl"),
							path.Join(createdTmpThemesDirectory, "default", "menu.tpl"),
							path.Join(createdTmpThemesDirectory, "default", "theme.css"),
						},
					}))
				})

			})

			Context("given the full theme path", func() {

				BeforeEach(func() {
					themeToOpen = path.Join(createdTmpThemesDirectory, "default")
				})

				It("should return the right theme", func() {
					Expect(returnedTheme).To(Equal(&Theme{
						Name: path.Join(createdTmpThemesDirectory, "default"),
						Path: path.Join(createdTmpThemesDirectory, "default"),
						Files: []string{
							path.Join(createdTmpThemesDirectory, "default", "index.tpl"),
							path.Join(createdTmpThemesDirectory, "default", "menu.tpl"),
							path.Join(createdTmpThemesDirectory, "default", "theme.css"),
						},
					}))
				})

			})

		})

		Context("when the theme does not exist", func() {

			BeforeEach(func() {
				themeToOpen = "theme-that-does-not-exist"
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(Equal("Theme not found"))
			})

			It("should not return a theme", func() {
				Expect(returnedTheme).To(BeNil())
			})

		})

		Context("when listing the theme files fails", func() {

			BeforeEach(func() {
				themeToOpen = "invalid-theme"
			})

			It("should return an error", func() {
				Expect(returnedError).NotTo(BeNil())
				Expect(returnedError.Error()).To(ContainSubstring("Failed to read theme directory:"))
			})

			It("should not return a theme", func() {
				Expect(returnedTheme).To(BeNil())
			})

		})

	})

})
