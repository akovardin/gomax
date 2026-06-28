package types

type Folder struct {
	SourceID   FlexInt                  `json:"sourceId"`
	Include    []FlexInt                `json:"include"`
	Options    map[string]interface{}   `json:"options"`
	UpdateTime FlexInt                  `json:"updateTime"`
	ID         FlexInt                  `json:"id"`
	Filters    []map[string]interface{} `json:"filters"`
	Title      string                   `json:"title"`
}

type FolderUpdate struct {
	FoldersOrder []FlexInt `json:"foldersOrder"`
	Folder       *Folder   `json:"folder"`
	FolderSync   FlexInt   `json:"folderSync"`
}

type FolderList struct {
	FoldersOrder            []FlexInt `json:"foldersOrder"`
	Folders                 []Folder  `json:"folders"`
	AllFilterExcludeFolders []FlexInt `json:"allFilterExcludeFolders"`
	FolderSync              FlexInt   `json:"folderSync"`
}
