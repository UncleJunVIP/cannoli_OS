package models

import (
	gcf "github.com/craterdog/go-collection-framework/v7"
)

type AppState struct {
	Config      *Config
	ScreenStack gcf.StackLike[string]
}
