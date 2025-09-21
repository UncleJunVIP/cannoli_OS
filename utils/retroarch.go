package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

var currentRetroArchProcess *os.Process

func LaunchRetroArchMenu() {
	logger := GetLoggerInstance()
	runRA([]string{"--menu", "-c", "retroarch.cfg"}, "RetroArch menu")

	logger.Debug("Sleeping for 2500ms before returning to menu")
	time.Sleep(2500 * time.Millisecond)
	logger.Debug("Sleep completed, returning to application")
}

func LaunchROM(gameName string, romPath string) {
	logger := GetLoggerInstance()
	logger.Debug(fmt.Sprintf("ROM path: %s", romPath))

	overlayClient := NewOverlayClient(gameName)
	err := overlayClient.Start()
	if err != nil {
		logger.Debug("Failed to start overlay", "error", err)
	} else {
		logger.Debug("IGM overlay started successfully")
	}

	corePath, err := determineCorePath(romPath)
	if err != nil {
		logger.Debug("Failed to determine core path", "error", err)
		return
	}

	process := runRA([]string{
		"-L", corePath,
		romPath,
	}, gameName)

	if process != nil {
		monitor := NewHotkeyMonitor()
		monitor.Start(overlayClient)

		process.Wait()

		monitor.Stop()
		logger.Debug("RetroArch exited, stopping overlay and monitor")
	}

	overlayClient.Stop()

	logger.Debug("Game session ended, returning to main menu")
}

func IsRetroArchRunning() bool {
	if currentRetroArchProcess == nil {
		return false
	}

	err := currentRetroArchProcess.Signal(syscall.Signal(0))
	return err == nil
}

func determineCorePath(romPath string) (string, error) {
	_, tag := ItemNameCleaner(filepath.Dir(romPath), false)

	core, exists := GetConfig().CoreMapping[tag]
	if !exists {
		return "", fmt.Errorf("could not determine core for ROM: %s", romPath)
	}

	ext, err := getCoreExtension()
	if err != nil {
		return "", err
	}

	coreFilename := core + "_libretro" + ext

	return filepath.Join(GetConfig().CoresDirectory, coreFilename), nil
}

func getCoreExtension() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return ".dll", nil
	case "darwin":
		return ".dylib", nil
	case "linux":
		return ".so", nil
	default:
		GetLoggerInstance().Error("Could not determine core extension for OS!")
		return "", fmt.Errorf("could not determine core extension for OS")
	}
}

func runRA(args []string, gameName string) *exec.Cmd {
	logger := GetLoggerInstance()
	logger.Debug(fmt.Sprintf("Starting %s", gameName))

	cmd := exec.Command("./retroarch", args...)
	cmd.Dir = config.RetroArchDirectory

	if os.Getenv("ENVIRONMENT") != "DEV" {
		cmd.Env = append(os.Environ(),
			"LD_LIBRARY_PATH=/mnt/SDCARD/System/RetroArch/lib:/usr/trimui/lib:"+os.Getenv("LD_LIBRARY_PATH"),
			"PATH=/usr/trimui/bin:"+os.Getenv("PATH"),
		)
	}

	// Get pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Debug("Failed to create stdout pipe: %v", err)
		return nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Debug("Failed to create stderr pipe!", "error", err)
		return nil
	}

	err = cmd.Start()
	if err != nil {
		logger.Debug("Failed to start RetroArch", "error", err)
		return nil
	}

	currentRetroArchProcess = cmd.Process
	logger.Debug(fmt.Sprintf("Started RetroArch with PID: %d", cmd.Process.Pid))

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			logger.Debug(fmt.Sprintf("[RA STDOUT] %s", scanner.Text()))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logger.Error(fmt.Sprintf("[RA STDERR] %s", scanner.Text()))
		}
	}()

	return cmd
}
