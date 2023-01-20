package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type FeatureServiceConfigDocumentInfo struct {
	Title    string `json:"Title"`
	Author   string `json:"Author"`
	Comments string `json:"Comments"`
	Subject  string `json:"Subject"`
	Category string `json:"Category"`
	Keywords string `json:"Keywords"`
}

type FeatureServiceConfig struct {
	CurrentVersion     float64 `json:"currentVersion"`
	ServiceDescription string  `json:"serviceDescription"`
	// Only map layers
	MapName                       string `json:"mapName"`
	HasVersionedData              bool   `json:"hasVersionedData"`
	HasArchivedData               bool   `json:"hasArchivedData"`
	SupportsDisconnectedEditing   bool   `json:"supportsDisconnectedEditing"`
	SupportsDatumTransformation   bool   `json:"supportsDatumTransformation"`
	SupportsRelationshipsResource bool   `json:"supportsRelationshipsResource"`
	SyncEnabled                   bool   `json:"syncEnabled"`
	SupportedQueryFormats         string `json:"supportedQueryFormats"`
	MaxRecordCount                int    `json:"maxRecordCount"`
	MaxRecordCountFactor          int    `json:"maxRecordCountFactor"`
	Capabilities                  string `json:"capabilities"`
	Description                   string `json:"description"`
	CopyrightText                 string `json:"copyrightText"`
	// Only map layers
	SupportsDynamicLayers bool                          `json:"supportsDynamicLayers"`
	SpatialReference      ServiceConfigSpatialReference `json:"spatialReference"`
	// Only map layers
	SingleFusedMapCache                         bool                `json:"singleFusedMapCache"`
	InitialExtent                               ServiceConfigExtent `json:"initialExtent"`
	FullExtent                                  ServiceConfigExtent `json:"fullExtent"`
	AllowGeometryUpdates                        bool                `json:"allowGeometryUpdates"`
	AllowTrueCurvesUpdates                      bool                `json:"allowTrueCurvesUpdates"`
	OnlyAllowTrueCurveUpdatesByTrueCurveClients bool                `json:"onlyAllowTrueCurveUpdatesByTrueCurveClients"`
	SupportsApplyEditsWithGlobalIds             bool                `json:"supportsApplyEditsWithGlobalIds"`
	SupportsTrueCurve                           bool                `json:"supportsTrueCurve"`
	Units                                       string              `json:"units"`
	// only map services
	SupportedImageFormatTypes string                             `json:"supportedImageFormatTypes"`
	DocumentInfo              FeatureServiceConfigDocumentInfo   `json:"documentInfo"`
	SupportsQueryDomains      bool                               `json:"supportsQueryDomains"`
	Layers                    []ServiceConfigLayer               `json:"layers"`
	Tables                    []ServiceConfigTable               `json:"tables"`
	Relationships             *[]any                             `json:"relationships"`
	EnableZDefaults           bool                               `json:"enableZDefaults"`
	ZDefault                  int                                `json:"zDefault"`
	AllowUpdateWithoutMValues bool                               `json:"allowUpdateWithoutMValues"`
	DatumTransformations      []ServiceConfigDatumTransformation `json:"datumTransformations"`
	ReferenceScale            float64                            `json:"referenceScale"`
}

func FetchFeatureServiceConfig(featureServiceURL *url.URL) (*FeatureServiceConfig, error) {
	rawURL := fmt.Sprintf("%s?f=json", featureServiceURL.String())

	page, err := fetch(rawURL)
	if err != nil {
		return nil, err
	}

	var featureServiceConfig FeatureServiceConfig
	err = json.Unmarshal(*page, &featureServiceConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal url: %s\n\n%s", rawURL, err))
	}

	return &featureServiceConfig, nil
}
