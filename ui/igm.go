package ui

import (
	"cannoliOS/models"

	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
)

type InGameMenu struct {
	Data     interface{}
	Position models.Position
	ROMPath  string
}

func (igm InGameMenu) Name() models.ScreenName {
	return models.InGameMenu
}

func (igm InGameMenu) Draw() (models.ScreenReturn, error) {
	menuItems := []gaba.MenuItem{
		{
			Text:     "Resume Game",
			Selected: false,
			Focused:  false,
			Metadata: "resume",
		},
		{
			Text:     "Save State",
			Selected: false,
			Focused:  false,
			Metadata: "save_state",
		},
		{
			Text:     "Load State",
			Selected: false,
			Focused:  false,
			Metadata: "load_state",
		},
		{
			Text:     "Reset Game",
			Selected: false,
			Focused:  false,
			Metadata: "reset",
		},
		{
			Text:     "Game Settings",
			Selected: false,
			Focused:  false,
			Metadata: "settings",
		},
		{
			Text:     "Quit",
			Selected: false,
			Focused:  false,
			Metadata: "quit",
		},
	}

	options := gaba.DefaultListOptions("In-Game Menu", menuItems)

	options.SelectedIndex = igm.Position.SelectedIndex
	options.VisibleStartIndex = igm.Position.SelectedPosition

	options.FooterHelpItems = []gaba.FooterHelpItem{
		{ButtonName: "B", HelpText: "Resume"},
		{ButtonName: "A", HelpText: "Select"},
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
