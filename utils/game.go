package utils

import (
	"fmt"
	"os/exec"
)

func LaunchROM(gameName string, romPath string) {
	logger := GetLoggerInstance()
	logger.Debug(fmt.Sprintf("ROM path: %s", romPath))

	igmPath := "./igm"
	cmd := exec.Command(igmPath, romPath)

	logger.Debug(fmt.Sprintf("Starting IGM for %s", gameName))

	err := cmd.Start()
	if err != nil {
		logger.Debug("Failed to start IGM", "error", err)
		return
	}

	logger.Debug(fmt.Sprintf("Started IGM with PID: %d", cmd.Process.Pid))

	err = cmd.Wait()
	if err != nil {
		logger.Debug("IGM process ended with error", "error", err)
	} else {
		logger.Debug("IGM process completed successfully")
	}

	logger.Debug("Game session ended, returning to main menu")
}
