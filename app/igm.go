package main

import (
	"cannoliOS/models"
	"cannoliOS/ui"
	"cannoliOS/utils"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/veandco/go-sdl2/sdl"
)

type OverlayResponse struct {
	Success bool   `json:"success"`
	Action  string `json:"action"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

const (
	shortPressMax = 2 * time.Second
	coolDownTime  = 1 * time.Second
)

var wg sync.WaitGroup

func main() {
	initLogging()

	utils.Logger.Println("Starting in-game overlay application...")

	gaba.InitSDL(gaba.GabagoolOptions{
		WindowTitle:    "In-Game Menu",
		ShowBackground: true,
		IsCannoli:      true,
	})

	menuButtonHandler(&wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	utils.Logger.Println("Shutting down overlay application...")
}

func initLogging() {
	logDir := "logs/igm"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory %s: %v", logDir, err)
		return
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("igm_%s.log", time.Now().Format("2006-01-02")))

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file %s: %v", logFile, err)
		return
	}

	utils.Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	utils.Logger.Printf("=== IGM Started at %s ===", time.Now().Format(time.RFC3339))
}

func menuButtonHandler(wg *sync.WaitGroup) {
	defer wg.Done()

	var pressTime time.Time
	var cooldownUntil time.Time

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			utils.Logger.Printf("Button event: %v", event)
			switch e := event.(type) {
			case *sdl.KeyboardEvent:

				if e.Keysym.Sym == sdl.K_1 {
					if e.State == sdl.RELEASED && !pressTime.IsZero() {
						log.Println("Short press detected, toggling menu...")
						cooldownUntil = time.Now().Add(coolDownTime)
						toggleMenu()
					} else if e.State == sdl.PRESSED {
						pressTime = time.Now()
					}
				}

			case *sdl.ControllerButtonEvent:
				if time.Now().Before(cooldownUntil) {
					continue
				}

				if gaba.Button(e.Button) == gaba.ButtonMenu {
					if e.State == sdl.RELEASED && !pressTime.IsZero() {
						log.Println("Short press detected, toggling menu...")
						cooldownUntil = time.Now().Add(coolDownTime)
						toggleMenu()
					} else if e.State == sdl.PRESSED {
						pressTime = time.Now()
					}
				}
			}
		}

		sdl.Delay(16) // ~60fps
	}
}

func toggleMenu() {
	sendRetroArchCommand("MUTE", "192.168.1.102", "55355")

	retroArchPID := getRetroArchPID()
	if retroArchPID > 0 {
		pauseRetroArch(retroArchPID)
	}

	time.Sleep(500 * time.Millisecond)

	gaba.ShowWindow()
	igm(utils.GetRomPath())

	if retroArchPID > 0 {
		resumeRetroArch(retroArchPID)
	}

	time.Sleep(300 * time.Millisecond)
	gaba.HideWindow()

	sendRetroArchCommand("MUTE", "192.168.1.102", "55355")
}

func getRetroArchPID() int {
	cmd := exec.Command("pgrep", "-f", "retroarch")
	output, err := cmd.Output()
	if err != nil {
		utils.Logger.Printf("Failed to find RetroArch process: %v", err)
		return 0
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		utils.Logger.Printf("No RetroArch process found")
		return 0
	}

	pids := strings.Split(pidStr, "\n")
	pid, err := strconv.Atoi(pids[0])
	if err != nil {
		utils.Logger.Printf("Failed to parse RetroArch PID: %v", err)
		return 0
	}

	return pid
}

func pauseRetroArch(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		utils.Logger.Printf("Failed to find RetroArch process %d: %v", pid, err)
		return
	}

	err = process.Signal(syscall.SIGSTOP)
	if err != nil {
		utils.Logger.Printf("Failed to pause RetroArch process %d: %v", pid, err)
		return
	}

	utils.Logger.Printf("Paused RetroArch process %d", pid)
}

func resumeRetroArch(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		utils.Logger.Printf("Failed to find RetroArch process %d: %v", pid, err)
		return
	}

	err = process.Signal(syscall.SIGCONT)
	if err != nil {
		utils.Logger.Printf("Failed to resume RetroArch process %d: %v", pid, err)
		return
	}

	utils.Logger.Printf("Resumed RetroArch process %d", pid)
}

func igm(romPath string) {
	utils.Logger.Printf("Showing in-game menu for ROM: %s", romPath)

	currentScreen := ui.InGameMenu{
		Data:     nil,
		Position: models.Position{},
		ROMPath:  romPath,
	}

	for {
		sr, err := currentScreen.Draw()
		if err != nil {
			utils.Logger.Printf("Error drawing in-game menu: %v", err)
		}

		switch sr.Code {
		case models.Back, models.Canceled:

		case models.Select:
			action := sr.Output.(string)
			utils.Logger.Printf("In-game menu action: %s", action)

			switch action {
			case "resume":
				return

			case "save_state":
				err := sendRetroArchCommand("SAVE_STATE", "192.168.1.102", "55355")
				if err != nil {
					utils.ShowMessage("Failed to save state", 3000)
				} else {
					utils.ShowMessage("Saved!", 3000)
				}
				return

			case "load_state":
				sendRetroArchCommand("LOAD_STATE", "192.168.1.102", "55355")
				return

			case "reset":
				sendRetroArchCommand("RESET", "192.168.1.102", "55355")
				return

			case "settings":
				retroArchPID := getRetroArchPID()
				resumeRetroArch(retroArchPID)
				time.Sleep(250 * time.Millisecond)
				sendRetroArchCommand("MENU_TOGGLE", "192.168.1.102", "55355")
				return

			case "quit":
				sendRetroArchCommand("QUIT", "192.168.1.102", "55355")
				return
			default:
				utils.Logger.Printf("Unhandled menu action: %s", action)
				continue
			}
		}
	}
}

func sendRetroArchCommand(command, host, port string) error {
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

	utils.Logger.Printf("Sent RetroArch UDP command: %s to %s:%s", command, host, port)
	return nil
}
