package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func ExecuteRetroArch(args []string, description string) {
	Logger.Printf("Starting %s", description)

	cmd := exec.Command("./retroarch", args...)
	cmd.Dir = "/mnt/SDCARD/RetroArch"

	Logger.Printf("Executing command: %s %v in directory: %s", cmd.Path, cmd.Args, cmd.Dir)

	cmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH=/mnt/SDCARD/RetroArch/lib:/usr/trimui/lib:"+os.Getenv("LD_LIBRARY_PATH"),
		"PATH=/usr/trimui/bin:"+os.Getenv("PATH"),
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	logAndPrintOutput(stdout.String(), stderr.String(), err, duration, description)
}

var currentRetroArchProcess *os.Process

func ExecuteRetroArchWithTracking(args []string, description string) *exec.Cmd {
	Logger.Printf("Starting %s", description)

	cmd := exec.Command("./retroarch", args...)
	cmd.Dir = "/mnt/SDCARD/RetroArch"

	cmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH=/mnt/SDCARD/RetroArch/lib:/usr/trimui/lib:"+os.Getenv("LD_LIBRARY_PATH"),
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

func logAndPrintOutput(stdoutStr, stderrStr string, err error, duration time.Duration, description string) {
	if stdoutStr != "" {
		Logger.Printf("Command stdout: %s", stdoutStr)
	}
	if stderrStr != "" {
		Logger.Printf("Command stderr: %s", stderrStr)
	}

	if err != nil {
		Logger.Printf("%s execution failed after %v: %v", description, duration, err)
		fmt.Println(err.Error())

		if stdoutStr != "" {
			fmt.Printf("stdout: %s\n", stdoutStr)
		}
		if stderrStr != "" {
			fmt.Printf("stderr: %s\n", stderrStr)
		}

		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			Logger.Printf("Process exited with code %d", exitError.ExitCode())
			fmt.Printf("Process exited with code %d\n", exitError.ExitCode())
		}
	} else {
		Logger.Printf("%s completed successfully after %v", description, duration)
		fmt.Println("Executable completed successfully")

		if stdoutStr != "" {
			fmt.Printf("stdout: %s\n", stdoutStr)
		}
	}
}
