package ui

import (
	"cannoliOS/models"
	"cannoliOS/state"
	"cannoliOS/utils"
	"path/filepath"
	"strings"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
	module "github.com/craterdog/go-collection-framework/v7"
)

type GameList struct {
	Directory      models.Directory
	SearchFilter   string
	DirectoryStack module.StackLike[models.Directory]
}

func (gl GameList) Name() models.ScreenName {
	return models.GameList
}

func (gl GameList) Draw() (models.ScreenReturn, error) {
	title := gl.Directory.DisplayName

	tagless, _ := utils.ItemNameCleaner(gl.Directory.DisplayName, true)
	if tagless != "" {
		title = tagless
	}

	fb := utils.NewFileBrowser()

	err := fb.CWD(gl.Directory.Path, false)
	if err != nil {
		// TODO fix this
	}

	var roms []models.Item
	roms = fb.Items

	if gl.SearchFilter != "" {
		title = "[Search: \"" + gl.SearchFilter + "\"]"
		//roms = utils.FilterList(roms, gl.SearchFilter)
	}

	var directoryEntries []gabagool.MenuItem
	var itemEntries []gabagool.MenuItem

	for _, item := range roms {
		if strings.HasPrefix(item.Filename, ".") { // Skip hidden files
			continue
		}

		itemName := strings.TrimSuffix(item.Filename, filepath.Ext(item.Filename))

		if item.IsMultiDiscDirectory || item.IsSelfContainedDirectory || !item.IsDirectory {
			imageFilename := strings.TrimSuffix(item.Filename, filepath.Ext(item.Filename)) + ".png"

			itemEntries = append(itemEntries, gabagool.MenuItem{
				Text:          itemName,
				Selected:      false,
				Focused:       false,
				Metadata:      item,
				ImageFilename: filepath.Join(gl.Directory.Path, ".media", imageFilename),
			})
		} else {
			itemName = "/" + itemName
			directoryEntries = append(directoryEntries, gabagool.MenuItem{
				Text:               itemName,
				Selected:           false,
				Focused:            false,
				Metadata:           item,
				NotMultiSelectable: true,
			})
		}
	}

	allEntries := append(directoryEntries, itemEntries...)

	options := gabagool.DefaultListOptions(title, allEntries)

	//selectedIndex, visibleStartIndex := state.GetCurrentMenuPosition()
	//options.SelectedIndex = selectedIndex
	//options.VisibleStartIndex = visibleStartIndex

	options.SmallTitle = false
	options.EmptyMessage = "No ROMs Found"
	options.EnableAction = true
	options.EnableMultiSelect = true
	options.FooterHelpItems = []gabagool.FooterHelpItem{
		{ButtonName: "B", HelpText: utils.GetString("back")},
		{ButtonName: "X", HelpText: utils.GetString("search")},
		{ButtonName: "Menu", HelpText: utils.GetString("help")},
	}

	appState := state.Get()

	if appState.Config.ShowArt {
		options.EnableImages = true
	}

	options.EnableHelp = true
	options.HelpTitle = "ROMs List Controls"
	options.HelpText = []string{
		"• X: Open Options",
		"• Select: Toggle Multi-Select",
		"• Start: Confirm Multi-Selection",
	}

	selection, err := gabagool.List(options)
	if err != nil {
		// TODO fix this
		return models.ScreenReturn{}, err
	}

	if selection.IsSome() && selection.Unwrap().ActionTriggered {
		//state.UpdateCurrentMenuPosition(selection.Unwrap().SelectedIndex, selection.Unwrap().VisiblePosition)
	} else if selection.IsSome() && !selection.Unwrap().ActionTriggered && selection.Unwrap().SelectedIndex != -1 {
		//state.UpdateCurrentMenuPosition(selection.Unwrap().SelectedIndex, selection.Unwrap().VisiblePosition)
		var selectedItems []models.Item
		rawSelection := selection.Unwrap().SelectedItems

		for _, item := range rawSelection {
			selectedItems = append(selectedItems, item.Metadata.(models.Item))
		}
		return models.ScreenReturn{
			Output:   selectedItems,
			Position: models.Position{},
			Code:     models.Select,
		}, nil
	} else if selection.IsSome() && selection.Unwrap().SelectedIndex == -1 {
		return models.ScreenReturn{
			Code: models.Back,
		}, nil
	}

	return models.ScreenReturn{
		Code: models.Canceled,
	}, nil
}
