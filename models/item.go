package models

type Item struct {
	DisplayName              string `json:"name"`
	Filename                 string `json:"filename"`
	Path                     string `json:"path"`
	IsDirectory              bool   `json:"is_directory"`
	IsSelfContainedDirectory bool   `json:"-"`
	IsMultiDiscDirectory     bool   `json:"-"`
	DirectoryFileCount       int    `json:"-"`
	FileSize                 string `json:"file_size"`
	LastModified             string `json:"last_modified"`
	Tag                      string `json:"tag"`
}

func (i Item) ToDirectory() Directory {
	return Directory{
		DisplayName: i.DisplayName,
		Tag:         i.Tag,
		Path:        i.Path,
	}
}
