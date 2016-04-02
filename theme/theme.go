package theme

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func GitClone(themeToDownload string, themesDirectory string, destinationDir string) error {
	var (
		err      error
		themeURL string
		byName   bool
	)
	if themeURL, byName, err = getThemeURL(themeToDownload); err != nil {
		return errors.New(fmt.Sprintf("An error occured while trying to resolve theme %v: %v", themeToDownload, err))
	}

	args := []string{"clone", themeURL}
	if destinationDir == "" {
		if byName {
			args = append(args, themeToDownload)
		}
	} else {
		args = append(args, destinationDir)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = themesDirectory
	stdErr := new(bytes.Buffer)
	cmd.Stderr = stdErr

	if err = cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("There was an error while executing git clone: %v", stdErr.String()))
	}

	return nil
}

func getThemeURL(themeToDownload string) (string, bool, error) {
	if match, _ := regexp.MatchString(`(https?:\/\/)|(git@)`, themeToDownload); match == true {
		return themeToDownload, false, nil
	}

	var (
		resp     *http.Response
		err      error
		maxTries int = 3
		i        int = 1
	)
	for {
		resp, err = http.Get("https://raw.githubusercontent.com/aubm/postmanerator-themes/master/.gitmodules")
		if err != nil {
			if i >= maxTries {
				return "", true, err
			}
			time.Sleep(5 * time.Second)
			i++
			continue
		}
		break
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		if scanner.Text() == fmt.Sprintf("\tpath = %v", themeToDownload) {
			scanner.Scan()
			urlLine := scanner.Text()
			urlLine = urlLine[7:len(urlLine)]
			return urlLine, true, nil
		}
	}

	return "", true, errors.New(fmt.Sprintf("Theme not found: %v", themeToDownload))
}

func Delete(themeToDelete string, themesDirectory string) error {
	return os.RemoveAll(fmt.Sprintf("%v/%v", themesDirectory, themeToDelete))
}

func ListThemes(out io.Writer, themesDirectory string) error {
	contents, err := readDir(themesDirectory)
	if err != nil {
		return errors.New("There was an error while reading the themes directory")
	}

	for _, dir := range contents {
		fmt.Fprintln(out, dir.Name())
	}

	return nil
}

func ListThemeFiles(themePath string) ([]string, error) {
	var files []string

	contents, err := readDir(themePath)
	if err != nil {
		return files, errors.New("There was an error while reading the theme contents")
	}

	for _, entry := range contents {
		if !entry.IsDir() {
			files = append(files, fmt.Sprintf("%v/%v", themePath, entry.Name()))
		}
	}

	return files, nil
}

func GetThemePath(themeId string, themesDirectory string) (string, error) {
	if _, err := os.Open(themeId); err == nil {
		return themeId, nil
	}

	localTheme := fmt.Sprintf("%v/%v", themesDirectory, themeId)
	if _, err := os.Open(localTheme); err == nil {
		return localTheme, nil
	}

	return "", errors.New("Theme not found")
}

func readDir(directory string) ([]os.FileInfo, error) {
	dirToRead, err := os.Open(directory)
	if err != nil {
		return []os.FileInfo{}, err
	}

	contents, err := dirToRead.Readdir(-1)
	if err != nil {
		return []os.FileInfo{}, err
	}

	return contents, nil
}
