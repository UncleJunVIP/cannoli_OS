package retroarch

import (
	"cannoliOS/utils"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Pause() {
	logger := utils.GetLoggerInstance()
	pid := getRetroArchPID()

	time.Sleep(250 * time.Millisecond)

	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Error("Failed to find RetroArch process", "error", err)
		return
	}

	err = process.Signal(syscall.SIGSTOP)
	if err != nil {
		logger.Error("Failed to pause RetroArch process", "error", err)
		return
	}

	logger.Debug("Paused RetroArch process", "pid", pid)
}

func Resume() {
	logger := utils.GetLoggerInstance()
	pid := getRetroArchPID()

	process, err := os.FindProcess(pid)
	if err != nil {
		logger.Error("Failed to find RetroArch process", "error", err, "pid", pid)
		return
	}

	err = process.Signal(syscall.SIGCONT)
	if err != nil {
		logger.Error("Failed to resume RetroArch process", "error", err, "pid", pid)
		return
	}

	logger.Debug("Resumed RetroArch process", "pid", pid)
}

func SendCommand(command, host, port string) error {
	logger := utils.GetLoggerInstance()

	addr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to connect to RetroArch UDP: %v", err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	_, err = conn.Write([]byte(command))
	if err != nil {
		return fmt.Errorf("failed to send UDP command: %v", err)
	}

	logger.Debug("Sent RetroArch UDP command", "command", command, "host", host, "port", port)
	return nil
}

func getRetroArchPID() int {
	logger := utils.GetLoggerInstance()

	cmd := exec.Command("pgrep", "-f", "retroarch")
	output, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to find RetroArch process", "error", err)
		return 0
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		logger.Debug("No RetroArch process found")
		return 0
	}

	pids := strings.Split(pidStr, "\n")
	pid, err := strconv.Atoi(pids[0])
	if err != nil {
		logger.Error("Failed to parse RetroArch PID", "error", err)
		return 0
	}

	return pid
}
