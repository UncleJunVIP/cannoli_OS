package utils

import (
	"cannoliOS/models"
	"cannoliOS/state"
	"fmt"
	"log/slog"
	"os"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
)

func Init() {
	var err error
	config, err := LoadConfig("config.json")

	if err != nil {
		GetLoggerInstance().Error("Failed to load config.json", "error", err)
		os.Exit(1)
	}

	gabagool.SetLogLevel(config.LogLevel)

	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		os.Exit(1)
	}

	state.Init(config)

	GetLoggerInstance().Info("=== Cannoli OS Started ===")
}

func GetLoggerInstance() *slog.Logger {
	return gabagool.GetLoggerInstance()
}

func GetConfig() *models.Config {
	return state.Get().Config
}
