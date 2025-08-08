package models

type ScreenName string

const (
	MainMenu  ScreenName = "MainMenu"
	GamesList ScreenName = "GameList"
	ToolsList ScreenName = "ToolsList"
	Settings  ScreenName = "Settings"
)

type Screen interface {
	Name() ScreenName
	Draw() (ScreenReturn, error)
}

type ScreenReturn struct {
	Output   interface{}
	Position Position
	Code     ExitCode
}

type Position struct {
	SelectedIndex    int
	SelectedPosition int
}
