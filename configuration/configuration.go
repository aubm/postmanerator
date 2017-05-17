package configuration

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
)

type Configuration struct {
	Out                    io.Writer
	ThemesRepository       string
	CollectionFile         string
	EnvironmentFile        string
	UsedTheme              string
	OutputFile             string
	Watch                  bool
	LocalName              string
	IgnoredRequestHeaders  StringsFlag
	IgnoredResponseHeaders StringsFlag
	ThemesDirectory        string
	Args                   []string
}

const (
	defaultThemesRepository = "https://raw.githubusercontent.com/aubm/postmanerator-themes/master/.gitmodules"
	postmaneratorPathEnv    = "POSTMANERATOR_PATH"
	homeEnv                 = "HOME"
	userProfileEnv          = "USERPROFILE"
)

var (
	InitErr error
	Config  = Configuration{
		Out:              os.Stdout,
		ThemesRepository: defaultThemesRepository,
	}
	errThemesNoDirectory = errors.New(`An error occurred while trying to determine which directory to use for themes.
As a workaround, you can define the POSTMANERATOR_PATH environment variable.
Please consult the documentation here https://github.com/aubm/postmanerator and feel free to submit an issue.`)
)

func init() {
	parseCommandFlags()
	parseCommandArgs()
	InitErr = parseThemesDir()
}

func parseCommandFlags() {
	flag.StringVar(&Config.CollectionFile, "collection", "", "the postman exported collection JSON file")
	flag.StringVar(&Config.EnvironmentFile, "environment", "", "the postman exported environment JSON file")
	flag.StringVar(&Config.UsedTheme, "theme", "default", "the theme to use")
	flag.StringVar(&Config.OutputFile, "output", "", "the output file, default is stdout")
	flag.BoolVar(&Config.Watch, "watch", false, "automatically regenerate the output when the theme changes")
	flag.StringVar(&Config.LocalName, "local-name", "", "the name of the local copy of the downloaded theme")
	flag.Var(&Config.IgnoredResponseHeaders, "ignored-response-headers", "a comma separated list of ignored response headers")
	flag.Var(&Config.IgnoredRequestHeaders, "ignored-request-headers", "a comma separated list of ignored request headers")
	flag.Parse()
}

func parseThemesDir() error {
	themesDir, err := getThemesDir()
	if err != nil {
		return err
	}

	if err := checkOrCreateThemesDir(themesDir); err != nil {
		return err
	}

	Config.ThemesDirectory = themesDir
	return nil
}

func getThemesDir() (string, error) {
	baseDir, err := getThemesBaseDir()
	themesDir := path.Join(baseDir, "themes")
	return themesDir, err
}

func getThemesBaseDir() (string, error) {
	themesDir := os.Getenv(postmaneratorPathEnv)
	if themesDir != "" {
		return themesDir, nil
	}

	usrHomeDir, err := getUserHomeDir()
	themesDir = path.Join(usrHomeDir, ".postmanerator")
	return themesDir, err
}

func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir, nil
	}

	if usrHomeDir := os.Getenv(homeEnv); usrHomeDir != "" {
		return usrHomeDir, nil
	}

	if usrHomeDir := os.Getenv(userProfileEnv); usrHomeDir != "" {
		return usrHomeDir, nil
	}

	return "", errThemesNoDirectory
}

func checkOrCreateThemesDir(themesDir string) error {
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(themesDir, 0777); err != nil {
			return fmt.Errorf("Failed to create themes directory: %v", err)
		}
	}
	return nil
}

func parseCommandArgs() {
	Config.Args = flag.Args()
}
