package commands

const (
	CmdDefault      = "cmd_default"
	CmdThemesList   = "cmd_themes_list"
	CmdThemesGet    = "cmd_themes_get"
	CmdThemesDelete = "cmd_themes_delete"
	CmdUnknown      = "cmd_unknown"
)

type Command interface {
	Is(name string) bool
	Do() error
}
