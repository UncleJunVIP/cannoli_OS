package ui

import (
	"cannoliOS/models"
	"cannoliOS/utils"

	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
)

type InGameMenu struct {
	Data     interface{}
	Position models.Position
	ROMPath  string
	GameName string
}

func (igm InGameMenu) Name() models.ScreenName {
	return models.InGameMenu
}

func (igm InGameMenu) Draw() (models.ScreenReturn, error) {
	menuItems := []gaba.MenuItem{
		{
			Text:     utils.GetString("resume"),
			Selected: false,
			Focused:  false,
			Metadata: "resume",
		},
		{
			Text:     utils.GetString("save_state"),
			Selected: false,
			Focused:  false,
			Metadata: "save_state",
		},
		{
			Text:     utils.GetString("load_state"),
			Selected: false,
			Focused:  false,
			Metadata: "load_state",
		},
		{
			Text:     utils.GetString("reset_game"),
			Selected: false,
			Focused:  false,
			Metadata: "reset",
		},
		{
			Text:     utils.GetString("settings"),
			Selected: false,
			Focused:  false,
			Metadata: "settings",
		},
		{
			Text:     utils.GetString("quit"),
			Selected: false,
			Focused:  false,
			Metadata: "quit",
		},
	}

	title := "In-Game Menu"

	if igm.GameName != "" {
		title = igm.GameName
	}

	options := gaba.DefaultListOptions(title, menuItems)

	options.SmallTitle = true
	options.SelectedIndex = igm.Position.SelectedIndex
	options.VisibleStartIndex = igm.Position.SelectedPosition

	options.FooterHelpItems = []gaba.FooterHelpItem{
		{ButtonName: "A", HelpText: utils.GetString("select")},
	}

	sel, err := gaba.List(options)
	if err != nil {
		return models.ScreenReturn{
			Code: models.Canceled,
		}, err
	}

	if sel.IsSome() {
		result := sel.Unwrap()

		if result.SelectedIndex == -1 {
			return models.ScreenReturn{
				Code: models.Back,
			}, nil
		}

		selectedAction := result.SelectedItem.Metadata.(string)

		return models.ScreenReturn{
			Output: selectedAction,
			Position: models.Position{
				SelectedIndex:    result.SelectedIndex,
				SelectedPosition: result.VisiblePosition,
			},
			Code: models.Select,
		}, nil
	}

	return models.ScreenReturn{
		Code: models.Canceled,
	}, nil
}
