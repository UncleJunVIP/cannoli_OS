package utils

import (
	"os"
	"os/exec"
	"syscall"
)

var currentRetroArchProcess *os.Process

func ExecuteRetroArchWithTracking(args []string, description string) *exec.Cmd {
	Logger.Printf("Starting %s", description)

	cmd := exec.Command("./retroarch", args...)
	cmd.Dir = "/mnt/SDCARD/System/RetroArch"

	cmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH=/mnt/SDCARD/System/RetroArch/lib:/usr/trimui/lib:"+os.Getenv("LD_LIBRARY_PATH"),
		"PATH=/usr/trimui/bin:"+os.Getenv("PATH"),
	)

	// Start but don't wait
	err := cmd.Start()
	if err != nil {
		Logger.Printf("Failed to start %s: %v", description, err)
		return nil
	}

	currentRetroArchProcess = cmd.Process
	Logger.Printf("Started %s with PID: %d", description, cmd.Process.Pid)

	return cmd
}

func IsRetroArchRunning() bool {
	if currentRetroArchProcess == nil {
		return false
	}

	err := currentRetroArchProcess.Signal(syscall.Signal(0))
	return err == nil
}
