package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type GPServiceConfig struct {
	CurrentVersion      float32  `json:"currentVersion"`
	CimVersion          string   `json:"cimVersion"`
	ServiceDescription  string   `json:"serviceDescription"`
	Tasks               []string `json:"tasks"`
	ExecutionType       string   `json:"executionType"`
	ResultMapServerName string   `json:"resultMapServerName"`
	MaximumRecords      int      `json:"maximumRecords"`
	Capabilities        string   `json:"capabilities"`
}

func FetchGPServiceConfig(gpServiceURL *url.URL) (*GPServiceConfig, error) {
	rawURL := fmt.Sprintf("%s?f=json", gpServiceURL.String())

	page, err := fetch(rawURL)
	if err != nil {
		return nil, err
	}

	var gpServiceConfig GPServiceConfig
	err = json.Unmarshal(*page, &gpServiceConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal url: %s\n\n%s", rawURL, err))
	}

	return &gpServiceConfig, nil
}
