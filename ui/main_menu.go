package ui

import (
	"cannoliOS/models"
	"cannoliOS/state"
	"cannoliOS/utils"
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"log"
)

type MainMenu struct {
	Data     interface{}
	Position models.Position
}

func (m MainMenu) Name() models.ScreenName {
	return models.MainMenu
}

func (m MainMenu) Draw() {
	var menuItems []gaba.MenuItem

	gameMenuItems, err := buildGameDirectoryMenuItems()
	if err != nil {
		// TODO fix this
	}

	menuItems = append(menuItems, gameMenuItems...)

	options := gaba.DefaultListOptions("cannoli_OS", menuItems)

	selectedIndex, visibleStartIndex := 0, 0 //TODO replace me with actual stack state
	options.SelectedIndex = selectedIndex
	options.VisibleStartIndex = visibleStartIndex

	options.EnableAction = true
	options.FooterHelpItems = []gaba.FooterHelpItem{
		{ButtonName: "B", HelpText: "Quit"},
		{ButtonName: "X", HelpText: "Settings"},
		{ButtonName: "A", HelpText: "Select"},
	}

	sel, err := gaba.List(options)
	if err != nil {
		// TODO do something
	}

	log.Printf("Selected: %v\n", sel)
}

func buildGameDirectoryMenuItems() ([]gaba.MenuItem, error) {
	fb := utils.NewFileBrowser()

	if err := fb.CWD(utils.GetRomPath(), state.Get().HideEmpty); err != nil {
		utils.ShowMessage("Error fetching ROM directories", 5000)
		return nil, err
	}

	var menuItems []gaba.MenuItem
	for _, item := range fb.Items {
		if item.IsDirectory {
			gameDirectory := item.ToDirectory()
			menuItems = append(menuItems, gaba.MenuItem{
				Text:     gameDirectory.DisplayName,
				Selected: false,
				Focused:  false,
				Metadata: gameDirectory,
			})
		}
	}

	return menuItems, nil
}
