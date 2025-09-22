package state

import (
	"cannoliOS/models"
	"sync"

	"go.uber.org/atomic"
)

var appState atomic.Pointer[models.AppState]
var onceAppState sync.Once

func Init(config *models.Config) {
	onceAppState.Do(func() {
		appState.Store(&models.AppState{
			Config: config,
		})
	})
}

func Get() *models.AppState {
	return appState.Load()
}

func Update(newAppState *models.AppState) {
	appState.Store(newAppState)
}
