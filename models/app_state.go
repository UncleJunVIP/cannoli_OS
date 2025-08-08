package models

import (
	gcf "github.com/craterdog/go-collection-framework/v7"
)

type AppState struct {
	HideEmpty   bool
	ScreenStack gcf.StackLike[string]
}
