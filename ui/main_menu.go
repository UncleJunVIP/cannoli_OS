package ui

import (
	"cannoliOS/models"
	"cannoliOS/state"
	"cannoliOS/utils"

	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
)

type MainMenu struct {
	Data     interface{}
	Position models.Position
}

func (m MainMenu) Name() models.ScreenName {
	return models.MainMenu
}

func (m MainMenu) Draw() (models.ScreenReturn, error) {
	var menuItems []gaba.MenuItem

	gameMenuItems := buildGameDirectoryMenuItems()

	menuItems = append(menuItems, gameMenuItems...)

	options := gaba.DefaultListOptions("cannoli_OS", menuItems)

	selectedIndex, visibleStartIndex := 0, 0 //TODO replace me with actual stack state
	options.SelectedIndex = selectedIndex
	options.VisibleStartIndex = visibleStartIndex
	options.DisableBackButton = true
	options.EnableMultiSelect = false
	options.EnableAction = true

	options.FooterHelpItems = []gaba.FooterHelpItem{
		{ButtonName: "X", HelpText: utils.GetString("tools")},
		{ButtonName: "A", HelpText: utils.GetString("select")},
	}

	sel, _ := gaba.List(options)

	if sel.IsSome() && sel.Unwrap().ActionTriggered {
		return models.ScreenReturn{
			Code: models.Action,
		}, nil
	} else if sel.IsSome() && !sel.Unwrap().ActionTriggered && sel.Unwrap().SelectedIndex != -1 {
		md := sel.Unwrap().SelectedItem.Metadata
		return models.ScreenReturn{
			Output: md,
			Position: models.Position{
				SelectedIndex:    sel.Unwrap().SelectedIndex,
				SelectedPosition: sel.Unwrap().VisiblePosition,
			},
			Code: models.Select,
		}, nil
	}

	return models.ScreenReturn{
		Code: models.Canceled,
	}, nil
}

func buildGameDirectoryMenuItems() []gaba.MenuItem {
	fb := utils.NewFileBrowser()

	if err := fb.CWD(utils.GetRomPath(), state.Get().Config.HideEmptyDirectories); err != nil {
		utils.GetLoggerInstance().Error("Failed to fetch ROM directories", "error", err)
		utils.ShowMessage("Error fetching ROM directories", 5000)
	}

	var menuItems []gaba.MenuItem
	for _, item := range fb.Items {
		if item.IsDirectory {
			gameDirectory := item.ToDirectory()
			tagless, _ := utils.ItemNameCleaner(gameDirectory.DisplayName, true)
			menuItems = append(menuItems, gaba.MenuItem{
				Text:     tagless,
				Selected: false,
				Focused:  false,
				Metadata: gameDirectory,
			})
		}
	}

	return menuItems
}
