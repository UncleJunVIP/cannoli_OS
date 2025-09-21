package models

import "log/slog"

type Config struct {
	RetroArchDirectory   string            `json:"retroarch_directory,omitempty"`
	CoresDirectory       string            `json:"cores_directory,omitempty"`
	CoreMapping          map[string]string `json:"core_mapping,omitempty"`
	ShowArt              bool              `json:"show_art,omitempty"`
	HideEmptyDirectories bool              `json:"hide_empty_directories,omitempty"`
	LogLevel             slog.Level        `json:"log_level,omitempty"`
}
