package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type OverlayClient struct {
	overlayProcess *os.Process
	isRunning      bool
	GameName       string
}

type OverlayCommand struct {
	Action  string `json:"action"`
	ROMPath string `json:"rom_path,omitempty"`
	Data    string `json:"data,omitempty"`
}

type OverlayResponse struct {
	Success bool   `json:"success"`
	Action  string `json:"action"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewOverlayClient(game string) *OverlayClient {
	return &OverlayClient{
		isRunning: false,
		GameName:  game,
	}
}

func (oc *OverlayClient) Start() error {
	if oc.isRunning {
		return nil
	}

	Logger.Println("Starting overlay application...")

	cmd := exec.Command("./igm", oc.GameName)
	cmd.Dir = "/mnt/SDCARD/System"

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start overlay application: %v", err)
	}

	oc.overlayProcess = cmd.Process
	oc.isRunning = true

	time.Sleep(500 * time.Millisecond)

	Logger.Printf("Overlay application started with PID: %d", cmd.Process.Pid)
	return nil
}

func (oc *OverlayClient) Stop() error {
	if !oc.isRunning || oc.overlayProcess == nil {
		return nil
	}

	Logger.Println("Stopping overlay application...")

	err := oc.overlayProcess.Signal(syscall.SIGTERM)
	if err != nil {
		Logger.Printf("Failed to send SIGTERM to overlay process: %v", err)
		err = oc.overlayProcess.Kill()
		if err != nil {
			Logger.Printf("Failed to kill overlay process: %v", err)
		}
	}

	_, waitErr := oc.overlayProcess.Wait()
	if waitErr != nil {
		Logger.Printf("Error waiting for overlay process to exit: %v", waitErr)
	}

	os.Remove("/tmp/cannoli_overlay.sock")

	oc.overlayProcess = nil
	oc.isRunning = false

	return err
}

func (oc *OverlayClient) ShowMenu() (*OverlayResponse, error) {
	return oc.sendCommand(OverlayCommand{
		Action: "show_menu",
	})
}

func (oc *OverlayClient) HideMenu() (*OverlayResponse, error) {
	return oc.sendCommand(OverlayCommand{
		Action: "hide_menu",
	})
}

func (oc *OverlayClient) SendRetroArchCommand(command string) (*OverlayResponse, error) {
	return oc.sendCommand(OverlayCommand{
		Action: "send_command",
		Data:   command,
	})
}

func (oc *OverlayClient) sendCommand(cmd OverlayCommand) (*OverlayResponse, error) {
	if !oc.isRunning {
		return nil, fmt.Errorf("overlay application not running")
	}

	conn, err := net.Dial("unix", "/tmp/cannoli_overlay.sock")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to overlay: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(cmd); err != nil {
		return nil, fmt.Errorf("failed to send command: %v", err)
	}

	decoder := json.NewDecoder(conn)
	var response OverlayResponse
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return &response, nil
}
