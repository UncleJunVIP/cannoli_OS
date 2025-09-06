package main

import (
	"bytes"
	"cannoliOS/models"
	"cannoliOS/ui"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	module "github.com/craterdog/go-collection-framework/v7"
)

var logger *log.Logger

func init() {
	initLogging()

	logger.Println("Initializing cannoli OS...")

	gaba.InitSDL(gaba.GabagoolOptions{
		WindowTitle:    "cannoli_OS",
		ShowBackground: true,
		IsCannoli:      true,
	})

	logger.Println("SDL initialization completed")
}

func initLogging() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("logs/cannoli_%s.log", timestamp)

	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}

	logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	logger.Printf("=== Cannoli OS Started at %s ===", time.Now().Format(time.RFC3339))
}

func main() {
	logger.Println("Starting main application loop")

	var currentScreen models.Screen

	currentScreen = ui.MainMenu{
		Data:     nil,
		Position: models.Position{},
	}

	logger.Printf("Initial screen set to: %s", currentScreen.Name())

	for {
		logger.Printf("Drawing screen: %s", currentScreen.Name())
		sr, err := currentScreen.Draw()

		if err != nil {
			logger.Printf("Error drawing screen %s: %v", currentScreen.Name(), err)
			continue
		}

		logger.Printf("Screen %s returned code: %v", currentScreen.Name(), sr.Code)

		switch currentScreen.Name() {
		case models.MainMenu:
			logger.Println("Processing MainMenu screen response")

			switch sr.Code {
			case models.Select:
				directory := sr.Output.(models.Directory)
				logger.Printf("Selected directory: %s (path: %s)", directory.DisplayName, directory.Path)

				if directory.DisplayName == "RetroArch" {
					launchRA()
					continue
				}

				currentScreen = ui.GameList{
					Directory:      directory,
					SearchFilter:   "",
					DirectoryStack: module.Stack[models.Directory](),
				}
				logger.Println("Switched to GameList screen")
			case models.Action:
				logger.Println("Action triggered in MainMenu")
			default:
				logger.Printf("Unhandled code in MainMenu: %v", sr.Code)
			}

		case models.GameList:
			logger.Println("Processing GameList screen response")
			gl := currentScreen.(ui.GameList)

			if sr.Code == models.Back {
				logger.Printf("Back action triggered, directory stack size: %d", gl.DirectoryStack.GetSize())

				if gl.DirectoryStack.GetSize() == 0 {
					logger.Println("Returning to MainMenu from GameList")
					currentScreen = ui.MainMenu{
						Data:     nil,
						Position: models.Position{},
					}
				} else {
					prev := gl.DirectoryStack.RemoveLast()
					logger.Printf("Navigating back to directory: %s", prev.DisplayName)
					currentScreen = ui.GameList{
						Directory:      prev,
						SearchFilter:   "",
						DirectoryStack: gl.DirectoryStack,
					}
				}
			} else if sr.Code == models.Select && sr.Output.([]models.Item)[0].IsDirectory { // TODO this needs to be cleaned
				selectedItem := sr.Output.([]models.Item)[0]
				logger.Printf("Selected directory item: %s", selectedItem.Filename)

				gl.DirectoryStack.AddValue(gl.Directory)
				currentScreen = ui.GameList{
					Directory:      selectedItem.ToDirectory(),
					SearchFilter:   "",
					DirectoryStack: gl.DirectoryStack,
				}
				logger.Printf("Navigated into directory: %s", selectedItem.Filename)
			} else if sr.Code == models.Select {
				selectedItems := sr.Output.([]models.Item)
				logger.Printf("Selected %d game item(s) for launch", len(selectedItems))
				for i, item := range selectedItems {
					logger.Printf("  Item %d: %s", i+1, item.Filename)
				}

				if len(selectedItems) > 0 {
					selectedItem := selectedItems[0]
					romPath := selectedItem.Path

					launchROM(romPath)
				}

				logger.Println("Returning to MainMenu after game launch")
				currentScreen = ui.MainMenu{
					Data:     nil,
					Position: models.Position{},
				}
			} else {

				logger.Printf("Unhandled code in GameList: %v", sr.Code)
			}
		default:
			logger.Printf("Unknown screen type: %s", currentScreen.Name())
		}
	}
}

func launchRA() {
	cmd := exec.Command("./retroarch", "--menu", "-c", "retroarch.cfg")

	cmd.Dir = "/mnt/SDCARD/RetroArch"

	logger.Printf("Executing command: %s %v in directory: %s", cmd.Path, cmd.Args, cmd.Dir)

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

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if stdoutStr != "" {
		logger.Printf("Command stdout: %s", stdoutStr)
	}
	if stderrStr != "" {
		logger.Printf("Command stderr: %s", stderrStr)
	}

	if err != nil {
		logger.Printf("RetroArch execution failed after %v: %v", duration, err)
		fmt.Println(err.Error())

		if stdoutStr != "" {
			fmt.Printf("stdout: %s\n", stdoutStr)
		}
		if stderrStr != "" {
			fmt.Printf("stderr: %s\n", stderrStr)
		}

		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			logger.Printf("Process exited with code %d", exitError.ExitCode())
			fmt.Printf("Process exited with code %d\n", exitError.ExitCode())
		}
	} else {
		logger.Printf("RetroArch completed successfully after %v", duration)
		fmt.Println("Executable completed successfully")

		if stdoutStr != "" {
			fmt.Printf("stdout: %s\n", stdoutStr)
		}
	}

	logger.Println("Sleeping for 1750ms before returning to menu")
	time.Sleep(1750 * time.Millisecond)
	logger.Println("Sleep completed, returning to application")
}

func launchROM(romPath string) {
	logger.Printf("Starting RetroArch execution with ROM: %s", romPath)

	var cmd *exec.Cmd

	cmd = exec.Command("./retroarch", "-L", "/mnt/SDCARD/RetroArch/cores/gambatte_libretro.so", romPath, "-c", "retroarch.cfg")

	cmd.Dir = "/mnt/SDCARD/RetroArch"

	logger.Printf("Executing command: %s %v in directory: %s", cmd.Path, cmd.Args, cmd.Dir)

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

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if stdoutStr != "" {
		logger.Printf("Command stdout: %s", stdoutStr)
	}
	if stderrStr != "" {
		logger.Printf("Command stderr: %s", stderrStr)
	}

	if err != nil {
		logger.Printf("RetroArch execution failed after %v: %v", duration, err)
		fmt.Println(err.Error())

		if stdoutStr != "" {
			fmt.Printf("stdout: %s\n", stdoutStr)
		}
		if stderrStr != "" {
			fmt.Printf("stderr: %s\n", stderrStr)
		}

		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			logger.Printf("Process exited with code %d", exitError.ExitCode())
			fmt.Printf("Process exited with code %d\n", exitError.ExitCode())
		}
	} else {
		logger.Printf("RetroArch completed successfully after %v", duration)
		fmt.Println("Executable completed successfully")

		if stdoutStr != "" {
			fmt.Printf("stdout: %s\n", stdoutStr)
		}
	}
}
