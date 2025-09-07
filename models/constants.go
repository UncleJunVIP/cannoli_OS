package models

type ExitCode int

const (
	Canceled ExitCode = -1
	Back     ExitCode = 0
	Select   ExitCode = 1
	Action   ExitCode = 2

	ResumeGame   = "resume_game"
	SaveState    = "save_state"
	LoadState    = "load_state"
	ResetGame    = "reset_game"
	GameSettings = "game_settings"
	ExitToMenu   = "exit_to_menu"
)
