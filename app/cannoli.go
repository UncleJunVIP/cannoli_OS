package main

import (
	"cannoliOS/models"
	"cannoliOS/ui"
	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
)

func init() {
	gaba.InitSDL(gaba.GabagoolOptions{
		WindowTitle:    "cannoli_OS",
		ShowBackground: true,
		IsCannoli:      true,
	})
}

func main() {
	mm := ui.MainMenu{
		Data:     nil,
		Position: models.Position{},
	}
	mm.Draw()
}
