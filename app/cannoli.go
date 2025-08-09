package main

import (
	"cannoliOS/models"
	"cannoliOS/ui"
	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	module "github.com/craterdog/go-collection-framework/v7"
)

func init() {
	gaba.InitSDL(gaba.GabagoolOptions{
		WindowTitle:    "cannoli_OS",
		ShowBackground: true,
		IsCannoli:      true,
	})
}

func main() {
	var currentScreen models.Screen

	currentScreen = ui.MainMenu{
		Data:     nil,
		Position: models.Position{},
	}

	for {
		sr, _ := currentScreen.Draw()

		switch currentScreen.Name() {
		case models.MainMenu:
			directory := sr.Output.(models.Directory)

			currentScreen = ui.GameList{
				Directory:      directory,
				SearchFilter:   "",
				DirectoryStack: module.Stack[models.Directory](),
			}

		case models.GameList:
			gl := currentScreen.(ui.GameList)
			if sr.Code == models.Back {
				if gl.DirectoryStack.GetSize() == 0 {
					currentScreen = ui.MainMenu{
						Data:     nil,
						Position: models.Position{},
					}
				} else {
					prev := gl.DirectoryStack.RemoveLast()
					currentScreen = ui.GameList{
						Directory:      prev,
						SearchFilter:   "",
						DirectoryStack: gl.DirectoryStack,
					}
				}
			} else if sr.Code == models.Select && sr.Output.([]models.Item)[0].IsDirectory { // TODO this needs to be cleaned
				gl.DirectoryStack.AddValue(gl.Directory)
				currentScreen = ui.GameList{
					Directory:      sr.Output.([]models.Item)[0].ToDirectory(),
					SearchFilter:   "",
					DirectoryStack: gl.DirectoryStack,
				}
			}
		}
	}
}
