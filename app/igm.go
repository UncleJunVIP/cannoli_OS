package main

import (
	"cannoliOS/models"
	"cannoliOS/retroarch"
	"cannoliOS/ui"
	"cannoliOS/utils"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	coolDownTime = 1 * time.Second
)

type menuAction string

const (
	Back       menuAction = ""
	ResumeGame            = "RESUME"
	SaveState             = "SAVE_STATE"
	LoadState             = "LOAD_STATE"
	Reset                 = "RESET"
	Settings              = "MENU_TOGGLE"
	Quit                  = "QUIT"
	Screenshot            = "SCREENSHOT"
)

var localIP = "127.0.0.1"

var romPath string
var gameName string

func init() {
	gaba.InitSDL(gaba.Options{
		WindowTitle:    "In-Game Menu",
		ShowBackground: true,
		IsCannoli:      true,
		LogFilename:    "igm.log",
		LogLevel:       slog.LevelDebug,
	})

	utils.LoadConfig()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: igm <rom file path>")
	}

	romPath = os.Args[1]
	gameName, _ = utils.ItemNameCleaner(filepath.Base(romPath), true)

	logger := utils.GetLoggerInstance()

	logger.Debug(fmt.Sprintf("Starting IGM for %s...", gameName))
	logger.Debug(fmt.Sprintf("ROM path: %s", romPath))

	gaba.HideWindow()

	retroArchCmd, err := retroarch.Launch(gameName, romPath)
	if err != nil {
		logger.Error("Failed to launch RetroArch", "error", err)
		return
	}

	retroArchExitChan := make(chan error, 1)

	go func() {
		err := retroArchCmd.Wait()
		retroArchExitChan <- err
	}()

	shutdownChan := make(chan bool, 1)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		logger.Debug("Received shutdown signal...")
		shutdownChan <- true
	}()

	go func() {
		err := <-retroArchExitChan
		if err != nil {
			logger.Debug("RetroArch process ended with error", "error", err)
		} else {
			logger.Debug("RetroArch process completed successfully")
		}
		shutdownChan <- true
	}()

	menuButtonHandler(shutdownChan)

	logger.Debug("Shutting down IGM...")

	if retroArchCmd.Process != nil {
		select {
		case <-retroArchExitChan:
			// Process already exited
		default:
			// Process might still be running, try to kill it
			retroArchCmd.Process.Kill()
		}
	}
}

func menuButtonHandler(shutdownChan <-chan bool) {
	var cooldownUntil time.Time

	for {
		select {
		case <-shutdownChan:
			return
		default:
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_1 {
					if e.State == sdl.PRESSED {
						log.Println("Button press detected, toggling menu...")
						cooldownUntil = time.Now().Add(coolDownTime)
						toggleMenu()
					}
				}

			case *sdl.ControllerButtonEvent:
				if time.Now().Before(cooldownUntil) {
					continue
				}

				if gaba.Button(e.Button) == gaba.ButtonMenu {
					if e.State == sdl.PRESSED {
						log.Println("Button press detected, toggling menu...")
						cooldownUntil = time.Now().Add(coolDownTime)
						toggleMenu()
					}
				}
			}
		}

		sdl.Delay(16)
	}
}

func toggleMenu() {
	logger := utils.GetLoggerInstance()

	retroarch.SendCommand(Screenshot, localIP, "55355")
	time.Sleep(175 * time.Millisecond)

	retroarch.Pause()

	time.Sleep(200 * time.Millisecond)

	gaba.ShowWindow()
	command, message := igm()

	logger.Debug(fmt.Sprintf("In-game menu command: %s", command))

	if command == Quit {
		gaba.ProcessMessage(fmt.Sprintf("%s %s...", message, gameName),
			gaba.ProcessMessageOptions{}, func() (interface{}, error) {
				retroarch.Terminate()
				time.Sleep(2000 * time.Millisecond)
				return nil, nil
			})
	} else {
		retroarch.Resume()
		if command != Back {
			retroarch.SendCommand(string(command), localIP, "55355")
		}
	}

	logger.Debug("Hiding IGM...")
	gaba.HideWindow()
}

func igm() (menuAction, string) {
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
			continue
		}

		switch sr.Code {
		case models.Back, models.Canceled:
			logger.Debug("Menu cancelled or back pressed")
			return Back, "" // Exit the menu without any command

		case models.Select:
			action := sr.Output.(string)
			logger.Debug("In-game menu action", "action", action)

			switch action {
			case "resume":
				return ResumeGame, ""

			case "save_state":
				return SaveState, utils.GetString("saving")

			case "load_state":
				return LoadState, utils.GetString("loading")

			case "reset":
				return Reset, utils.GetString("resetting")

			case "settings":
				return Settings, ""

			case "quit":
				return Quit, utils.GetString("quitting")
			default:
				logger.Debug("Unhandled menu action", "action", action)
				continue
			}
		default:
			logger.Debug("Unhandled screen response code", "code", sr.Code)
		}
	}
}
