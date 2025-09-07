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

	gameMenuItems, err := buildGameDirectoryMenuItems()
	if err != nil {
		// TODO fix this
	}

	menuItems = append(menuItems, gameMenuItems...)

	menuItems = append(menuItems, gaba.MenuItem{
		Text:     "RetroArch",
		Metadata: models.Directory{DisplayName: "RetroArch"},
	})

	options := gaba.DefaultListOptions("cannoli_OS", menuItems)

	selectedIndex, visibleStartIndex := 0, 0 //TODO replace me with actual stack state
	options.SelectedIndex = selectedIndex
	options.VisibleStartIndex = visibleStartIndex
	options.DisableBackButton = true
	options.EnableMultiSelect = false

	options.EnableAction = true
	options.FooterHelpItems = []gaba.FooterHelpItem{
		{ButtonName: "X", HelpText: "Settings"},
		{ButtonName: "A", HelpText: "Select"},
	}

	sel, err := gaba.List(options)
	if err != nil {
		// TODO do something
	}

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

func buildGameDirectoryMenuItems() ([]gaba.MenuItem, error) {
	fb := utils.NewFileBrowser()

	if err := fb.CWD(utils.GetRomPath(), state.Get().Config.HideEmptyDirectories); err != nil {
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
