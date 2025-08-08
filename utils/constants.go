package utils

import (
	"os"
	"path/filepath"
	"regexp"
)

const (
	systemPath      = "/System"
	biosPath        = "/BIOS"
	romsPath        = "/Roms"
	cheatsPath      = "/Cheats"
	collectionsPath = "/Collections"
	archivesPath    = "/Archives"
	savesPath       = "/Saves"
	saveStatePath   = "/Save States"
	settingsPath    = "/Settings"
)

var OrderedFolderRegex = regexp.MustCompile(`\d+\)\s`)
var TagRegex = regexp.MustCompile(`\((.*?)\)`)

func GetRoot() string {
	return os.Getenv("HOME")
}

func GetSystemPath() string {
	return filepath.Join(GetRoot(), systemPath)
}

func GetBiosPath() string {
	return filepath.Join(GetRoot(), biosPath)
}

func GetRomPath() string {
	return filepath.Join(GetRoot(), romsPath)
}

func GetCollectionPath() string {
	return filepath.Join(GetRoot(), collectionsPath)
}

func GetArchivePath() string {
	return filepath.Join(GetRoot(), archivesPath)
}

func GetSavePath() string {
	return filepath.Join(GetRoot(), savesPath)
}

func GetSaveStatePath() string {
	return filepath.Join(GetRoot(), saveStatePath)
}

func GetCheatsPath() string {
	return filepath.Join(GetRoot(), cheatsPath)
}

func GetSettingsPath() string {
	return filepath.Join(GetRoot(), settingsPath)
}
