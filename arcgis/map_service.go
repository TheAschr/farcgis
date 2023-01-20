package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type MapServiceConfigDocumentInfo struct {
	Title                string `json:"Title"`
	Author               string `json:"Author"`
	Comments             string `json:"Comments"`
	Subject              string `json:"Subject"`
	Category             string `json:"Category"`
	AntialiasingMode     string `json:"AntialiasingMode"`
	TextAntialiasingMode string `json:"TextAntialiasingMode"`
	Keywords             string `json:"Keywords"`
}

type MapServiceConfig struct {
	CurrentVersion              float64                       `json:"currentVersion"`
	ServiceDescription          string                        `json:"serviceDescription"`
	MapName                     string                        `json:"mapName"`
	Description                 string                        `json:"description"`
	CopyrightText               string                        `json:"copyrightText"`
	SupportsDynamicLayers       bool                          `json:"supportsDynamicLayers"`
	Layers                      []ServiceConfigLayer          `json:"layers"`
	Tables                      []ServiceConfigTable          `json:"tables"`
	SpatialReference            ServiceConfigSpatialReference `json:"spatialReference"`
	SingleFusedMapCache         bool                          `json:"singleFusedMapCache"`
	InitialExtent               ServiceConfigExtent           `json:"initialExtent"`
	FullExtent                  ServiceConfigExtent           `json:"fullExtent"`
	MinScale                    float64                       `json:"minScale"`
	MaxScale                    float64                       `json:"maxScale"`
	Units                       string                        `json:"units"`
	SupportedImageFormatTypes   string                        `json:"supportedImageFormatTypes"`
	DocumentInfo                MapServiceConfigDocumentInfo  `json:"documentInfo"`
	Capabilities                string                        `json:"capabilities"`
	SupportedQueryFormats       string                        `json:"supportedQueryFormats"`
	HasVersionedData            bool                          `json:"hasVersionedData"`
	ExportTilesAllowed          bool                          `json:"exportTilesAllowed"`
	ReferenceScale              float64                       `json:"referenceScale"`
	SupportsDatumTransformation bool                          `json:"supportsDatumTransformation"`
	MaxRecordCount              int                           `json:"maxRecordCount"`
	MaxImageHeight              int                           `json:"maxImageHeight"`
	MaxImageWidth               int                           `json:"maxImageWidth"`
	SupportedExtensions         string                        `json:"supportedExtensions"`
}

func FetchMapServiceConfig(mapServiceURL *url.URL) (*MapServiceConfig, error) {
	rawURL := fmt.Sprintf("%s?f=json", mapServiceURL.String())

	page, err := fetch(rawURL)
	if err != nil {
		return nil, err
	}

	var mapServiceConfig MapServiceConfig
	err = json.Unmarshal(*page, &mapServiceConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal url: %s\n\n%s", rawURL, err))
	}

	return &mapServiceConfig, nil
}
