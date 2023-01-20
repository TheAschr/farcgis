package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type GeometryServiceConfig struct {
	ServiceDescription string `json:"serviceDescription"`
}

func FetchGeometryServiceConfig(geometryServiceURL *url.URL) (*GeometryServiceConfig, error) {
	rawURL := fmt.Sprintf("%s?f=json", geometryServiceURL.String())

	page, err := fetch(rawURL)
	if err != nil {
		return nil, err
	}

	var geometryServiceConfig GeometryServiceConfig
	err = json.Unmarshal(*page, &geometryServiceConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal url: %s\n\n%s", rawURL, err))
	}

	return &geometryServiceConfig, nil
}
