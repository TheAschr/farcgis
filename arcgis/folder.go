package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type FolderConfigService struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FolderConfig struct {
	CurrentVersion float32               `json:"currentVersion"`
	Folders        []string              `json:"folders"`
	Services       []FolderConfigService `json:"services"`
}

func FetchFolderConfig(folderURL *url.URL) (*FolderConfig, error) {
	rawURL := fmt.Sprintf("%s?f=json", folderURL.String())

	page, err := fetch(rawURL)
	if err != nil {
		return nil, err
	}

	var folderConfig FolderConfig
	err = json.Unmarshal(*page, &folderConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal url: %s\n\n%s", rawURL, err))
	}

	return &folderConfig, nil
}
