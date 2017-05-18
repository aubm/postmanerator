package themes

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/aubm/postmanerator/configuration"
	"github.com/aubm/postmanerator/utils"
)

var (
	ErrThemeNotFound = errors.New("Theme not found")
	gitUrlRegexp     = regexp.MustCompile(`(https?:\/\/)|(git@)`)
)

type Manager struct {
	Config *configuration.Configuration `inject:""`
	Cloner interface {
		Clone(args []string, options utils.CloneOptions) error
	} `inject:""`
}

func (m *Manager) Delete(theme string) error {
	theme = fmt.Sprintf("%s/%s", m.Config.ThemesDirectory, theme)
	return os.RemoveAll(theme)
}

func (m *Manager) List() ([]string, error) {
	themeDirs, err := m.readDir(m.Config.ThemesDirectory)
	if err != nil {
		return nil, fmt.Errorf("There was an error while reading the themes directory: %v", err)
	}

	themeList := make([]string, 0)
	for _, themeDir := range themeDirs {
		themeList = append(themeList, themeDir.Name())
	}

	return themeList, nil
}

func (m *Manager) Download(theme string) (err error) {
	var localName string

	if !m.isGitUrl(theme) {
		localName = theme
		theme, err = m.getThemeURL(theme)
		if err != nil {
			return
		}
	}

	return m.clone(theme, localName)
}

func (m *Manager) isGitUrl(theme string) bool {
	return gitUrlRegexp.MatchString(theme)
}

func (m *Manager) getThemeURL(themeName string) (string, error) {
	r, err := m.getThemesListReader()
	if err != nil {
		return "", fmt.Errorf("failed to download the theme list: %v", err)
	}
	defer r.Close()

	return m.searchForThemeUrlInThemeList(r, themeName)
}

func (m *Manager) getThemesListReader() (io.ReadCloser, error) {
	var resp *http.Response
	var err error
	nbMaxTries := 3
	sleepTimeBetweenEachTry := 3 * time.Second
	for i := 0; i < nbMaxTries; i++ {
		resp, err = http.Get(m.Config.ThemesRepository)
		if err == nil {
			return resp.Body, nil
		}
		if i < nbMaxTries {
			time.Sleep(sleepTimeBetweenEachTry)
		}
	}
	return nil, err
}

func (m *Manager) searchForThemeUrlInThemeList(themeList io.Reader, themeNameToSearch string) (string, error) {
	scanner := bufio.NewScanner(themeList)
	for scanner.Scan() {
		if scanner.Text() == fmt.Sprintf("\tpath = %v", themeNameToSearch) {
			scanner.Scan()
			urlLine := scanner.Text()
			urlLine = urlLine[7:len(urlLine)]
			return urlLine, nil
		}
	}
	return "", ErrThemeNotFound
}

func (m *Manager) clone(url string, localName string) error {
	args := []string{url}

	if m.Config.ThemeLocalName != "" {
		localName = m.Config.ThemeLocalName
	}

	if localName != "" {
		args = append(args, localName)
	}

	options := utils.CloneOptions{TargetDirectory: m.Config.ThemesDirectory}
	return m.Cloner.Clone(args, options)
}

func (m *Manager) Open(themeName string) (*Theme, error) {
	themePath, err := m.getThemePath(themeName)
	if err != nil {
		return nil, err
	}

	theme := &Theme{Name: themeName, Path: themePath}
	theme.Files, err = m.listThemeFiles(themePath)
	if err != nil {
		return nil, err
	}

	return theme, nil
}

func (m *Manager) getThemePath(theme string) (string, error) {
	if ok := m.directoryExists(theme); ok {
		return theme, nil
	}

	theme = fmt.Sprintf("%v/%v", m.Config.ThemesDirectory, theme)
	if ok := m.directoryExists(theme); ok {
		return theme, nil
	}

	return "", ErrThemeNotFound
}

func (m *Manager) directoryExists(dir string) bool {
	if f, err := os.Open(dir); err == nil {
		f.Close()
		return true
	}
	return false
}

func (m *Manager) listThemeFiles(themePath string) ([]string, error) {
	contents, err := m.readDir(themePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read theme directory: %v", err)
	}

	themeFiles := make([]string, 0)
	for _, entry := range contents {
		if !entry.IsDir() {
			file := path.Join(themePath, entry.Name())
			themeFiles = append(themeFiles, file)
		}
	}

	return themeFiles, nil
}

func (m *Manager) readDir(directory string) ([]os.FileInfo, error) {
	dirToRead, err := os.Open(directory)
	if err != nil {
		return nil, err
	}

	contents, err := dirToRead.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
