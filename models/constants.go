package models

type ExitCode int

const (
	Canceled ExitCode = -1
	Back     ExitCode = 0
	Select   ExitCode = 1
	Action   ExitCode = 2
)
