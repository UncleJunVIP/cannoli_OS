package main

import (
	"cannoliOS/models"
	"cannoliOS/ui"
	"cannoliOS/utils"
	"time"

	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	module "github.com/craterdog/go-collection-framework/v7"
)

func init() {
	utils.Logger.Println("Initializing cannoli OS...")

	gaba.InitSDL(gaba.GabagoolOptions{
		WindowTitle:    "cannoli_OS",
		ShowBackground: true,
		IsCannoli:      true,
	})

	utils.Logger.Println("SDL initialization completed")
}

func main() {
	logger := utils.Logger
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

					gaba.HideWindow()
					launchROM(selectedItem.DisplayName, romPath)
					gaba.ShowWindow()
				}

				logger.Println("Returning to MainMenu after game launch")
				currentScreen = ui.MainMenu{
					Data:     nil,
					Position: models.Position{},
				}
				time.Sleep(1250 * time.Millisecond)
			} else {

				logger.Printf("Unhandled code in GameList: %v", sr.Code)
			}

		default:
			logger.Printf("Unknown screen type: %s", currentScreen.Name())
		}
	}
}

func launchRA() {
	logger := utils.Logger
	utils.ExecuteRetroArchWithTracking([]string{"--menu", "-c", "retroarch.cfg"}, "RetroArch menu")

	logger.Println("Sleeping for 2500ms before returning to menu")
	time.Sleep(2500 * time.Millisecond)
	logger.Println("Sleep completed, returning to application")
}

func launchROM(name string, romPath string) {
	logger := utils.Logger
	logger.Printf("ROM path: %s", romPath)

	overlayClient := utils.NewOverlayClient(name)
	err := overlayClient.Start()
	if err != nil {
		logger.Printf("Failed to start overlay: %v", err)
	} else {
		logger.Println("IGM overlay started successfully")
	}

	process := utils.ExecuteRetroArchWithTracking([]string{
		"-L", "/mnt/SDCARD/System/RetroArch/cores/gambatte_libretro.so",
		romPath,
		"-c", "retroarch.cfg",
	}, "RetroArch with ROM")

	if process != nil {
		monitor := utils.NewHotkeyMonitor()
		monitor.Start(overlayClient)

		process.Wait()

		monitor.Stop()
		logger.Println("RetroArch exited, stopping overlay and monitor")
	}

	overlayClient.Stop()

	logger.Println("Game session ended, returning to main menu")
}
