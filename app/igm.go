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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	coolDownTime = 1 * time.Second
)

var localIP = getIPFromInterface("wlan0")

var gameName string

var wg sync.WaitGroup

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: igm <game_name>")
	}

	gameName = os.Args[1]

	logger := utils.GetLoggerInstance()

	logger.Debug(fmt.Sprintf("Starting in-game overlay application for %s...", gameName))

	gaba.InitSDL(gaba.Options{
		WindowTitle:    "In-Game Menu",
		ShowBackground: true,
		IsCannoli:      true,
	})

	menuButtonHandler(&wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Debug("Shutting down overlay application...")
}

func menuButtonHandler(wg *sync.WaitGroup) {
	defer wg.Done()

	var pressTime time.Time
	var cooldownUntil time.Time

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
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
	logger := utils.GetLoggerInstance()

	retroArchPID := getRetroArchPID()
	if retroArchPID > 0 {
		pauseRetroArch(retroArchPID)
	}

	time.Sleep(200 * time.Millisecond)

	gaba.ShowWindow()
	command, message := igm()

	logger.Debug("In-game menu command: %s", command)

	if command != "" {
		if message != "" {
			gaba.ProcessMessage(fmt.Sprintf("%s...", message), gaba.ProcessMessageOptions{}, func() (interface{}, error) {
				sendRetroArchCommand(command, localIP, "55355", true)
				return nil, nil
			})
		} else {
			sendRetroArchCommand(command, localIP, "55355", true)
		}
	} else {
		resumeRetroArch(retroArchPID)
		time.Sleep(250 * time.Millisecond)
	}

	gaba.HideWindow()
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

func pauseRetroArch(pid int) {
	logger := utils.GetLoggerInstance()

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

func resumeRetroArch(pid int) {
	logger := utils.GetLoggerInstance()

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

func igm() (string, string) {
	logger := utils.GetLoggerInstance()

	logger.Debug("Showing in-game menu for ROM", "game_name", gameName)

	currentScreen := ui.InGameMenu{
		Data:     nil,
		Position: models.Position{},
		GameName: gameName,
	}

	for {
		sr, err := currentScreen.Draw()
		if err != nil {
			logger.Error("Error drawing in-game menu", "error", err)
		}

		switch sr.Code {
		case models.Back, models.Canceled:

		case models.Select:
			action := sr.Output.(string)
			logger.Debug("In-game menu action", "action", action)

			switch action {
			case "resume":
				return "", ""

			case "save_state":
				return "SAVE_STATE", utils.GetString("saving")

			case "load_state":
				return "LOAD_STATE", utils.GetString("loading")

			case "reset":
				return "RESET", utils.GetString("resetting")

			case "settings":
				return "MENU_TOGGLE", ""

			case "quit":
				return "QUIT", utils.GetString("quitting")
			default:
				logger.Debug("Unhandled menu action", "action", action)
				continue
			}
		}
	}
}

func sendRetroArchCommand(command, host, port string, resume bool) error {
	logger := utils.GetLoggerInstance()

	retroArchPID := getRetroArchPID()

	if resume {
		time.Sleep(750 * time.Millisecond)
		resumeRetroArch(retroArchPID)
		time.Sleep(250 * time.Millisecond)
	}

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

func getIPFromInterface(interfaceName string) string {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
