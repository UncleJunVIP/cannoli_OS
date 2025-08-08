package state

import (
	"cannoliOS/models"
	"go.uber.org/atomic"
	"sync"
)

var appState atomic.Pointer[models.AppState]
var onceAppState sync.Once

func Get() *models.AppState {
	onceAppState.Do(func() {
		appState.Store(&models.AppState{
			HideEmpty: true,
		})
	})
	return appState.Load()
}

func Update(newAppState *models.AppState) {
	appState.Store(newAppState)
}
