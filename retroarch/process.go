package retroarch

import (
	"bufio"
	"cannoliOS/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

func Launch(gameName string, romPath string) (*exec.Cmd, error) {
	logger := utils.GetLoggerInstance()

	corePath, err := determineCorePath(romPath)
	if err != nil {
		return nil, fmt.Errorf("failed to determine core path: %v", err)
	}

	args := []string{
		"-L", corePath,
		romPath,
	}

	cmd := exec.Command("./retroarch", args...)

	cmd.Dir = utils.GetConfig().RetroArchDirectory

	if os.Getenv("ENVIRONMENT") != "DEV" {
		cmd.Env = append(os.Environ(),
			"LD_LIBRARY_PATH=/mnt/SDCARD/System/RetroArch/lib:/usr/trimui/lib:"+os.Getenv("LD_LIBRARY_PATH"),
			"PATH=/usr/trimui/bin:"+os.Getenv("PATH"),
		)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start RetroArch: %v", err)
	}

	logger.Debug(fmt.Sprintf("Started RetroArch with PID: %d for game: %s", cmd.Process.Pid, gameName))

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			//logger.Debug(fmt.Sprintf("[RA STDOUT] %s", scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			//logger.Error(fmt.Sprintf("[RA STDOUT] Scanner error: %v", err))
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			//logger.Error(fmt.Sprintf("[RA STDERR] %s", scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			//logger.Error(fmt.Sprintf("[RA STDERR] Scanner error: %v", err))
		}
	}()

	return cmd, nil
}

func Terminate() {
	logger := utils.GetLoggerInstance()
	pid := getRetroArchPID()

	if pid == 0 {
		logger.Debug("No RetroArch process found to terminate")
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Error("Failed to find RetroArch process", "error", err, "pid", pid)
		return
	}

	err = process.Signal(syscall.SIGKILL)
	if err != nil {
		logger.Error("Failed to terminate RetroArch process", "error", err, "pid", pid)
		return
	}

	logger.Debug("Sent SIGKILL to RetroArch process", "pid", pid)
}

func determineCorePath(romPath string) (string, error) {
	_, tag := utils.ItemNameCleaner(filepath.Dir(romPath), false)

	core, exists := utils.GetConfig().CoreMapping[tag]
	if !exists {
		return "", fmt.Errorf("could not determine core for ROM: %s", romPath)
	}

	ext, err := getCoreExtension()
	if err != nil {
		return "", err
	}

	coreFilename := core + "_libretro" + ext

	return filepath.Join(utils.GetConfig().CoresDirectory, coreFilename), nil
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
		utils.GetLoggerInstance().Error("Could not determine core extension for OS!")
		return "", fmt.Errorf("could not determine core extension for OS")
	}
}
