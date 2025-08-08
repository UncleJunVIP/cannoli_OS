package app

import (
	_ "github.com/UncleJunVIP/certifiable"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
)

func init() {
	gaba.InitSDL(gaba.GabagoolOptions{
		WindowTitle:    "cannoliOS",
		ShowBackground: true,
	})
}
