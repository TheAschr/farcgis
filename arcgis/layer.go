package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type LayerConfigDrawingInfoRendererSymbolOutline struct {
	Color [4]int  `json:"color"`
	Width float32 `json:"width"`
}

type LayerConfigDrawingInfoRendererSymbol struct {
	Type    string                                      `json:"type"`
	Style   string                                      `json:"style"`
	Color   [4]int                                      `json:"color"`
	Size    int                                         `json:"size"`
	Angle   int                                         `json:"angle"`
	XOffset int                                         `json:"xoffset"`
	YOffset int                                         `json:"yoffset"`
	Outline LayerConfigDrawingInfoRendererSymbolOutline `json:"outline"`
}

type LayerConfigDrawingInfoRenderer struct {
	Type        string                               `json:"type"`
	Symbol      LayerConfigDrawingInfoRendererSymbol `json:"symbol"`
	Label       string                               `json:"label"`
	Description string                               `json:"description"`
}

type LayerConfigDrawingInfo struct {
	Renderer     LayerConfigDrawingInfoRenderer `json:"renderer"`
	Transparency int                            `json:"transparency"`
	LabelingInfo *any                           `json:"labelingInfo"`
}

type LayerConfigSpatialReference struct {
	WKID       int `json:"wkid"`
	LatestWKID int `json:"lastestWkid"`
}

type LayerConfigExtentSpatialReference struct {
	WKID       int `json:"wkid"`
	LatestWKID int `json:"lastestWkid"`
}

type LayerConfigExtent struct {
	XMin             interface{}                       `json:"xmin"` // Can be float or "NaN"
	YMin             interface{}                       `json:"ymin"` // Can be float or "NaN"
	XMax             interface{}                       `json:"xmax"` // Can be float or "NaN"
	YMax             interface{}                       `json:"ymax"` // Can be float or "NaN"
	SpatialReference LayerConfigExtentSpatialReference `json:"spatialReference"`
}

type LayerConfigFieldDomainCodedValue struct {
	Name string      `json:"name"`
	Code interface{} `json:"code"` // Can be string or int
}

type LayerConfigFieldDomain struct {
	Type        string                             `json:"type"`
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	CodedValues []LayerConfigFieldDomainCodedValue `json:"codedValues"`
	MergePolicy string                             `json:"mergePolicy"`
	SplitPolicy string                             `json:"splitPolicy"`
}

type LayerConfigField struct {
	Name   string                  `json:"name"`
	Type   string                  `json:"type"`
	Alias  string                  `json:"alias"`
	Length *int                    `json:"length,omitempty"`
	Domain *LayerConfigFieldDomain `json:"domain"`
}

type LayerConfigGeometryField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Alias string `json:"alias"`
}

type LayerConfigIndex struct {
	Name        string `json:"name"`
	Fields      string `json:"fields"`
	IsAscending bool   `json:"isAscending"`
	IsUnique    bool   `json:"isUnique"`
	Description string `json:"description"`
}

type LayerConfigOwnershipBasedAccessControlForFeatures struct {
	AllowOthersToQuery bool `json:"allowOthersToQuery"`
}

type LayerConfigAdvancedQueryCapabilities struct {
	UseStandardizedQueries       bool `json:"useStandardizedQueries"`
	SupportsStatistics           bool `json:"supportsStatistics"`
	SupportsHavingClause         bool `json:"supportsHavingClause"`
	SupportsCountDistinct        bool `json:"supportsCountDistinct"`
	SupportsOrderBy              bool `json:"supportsOrderBy"`
	SupportsDistinct             bool `json:"supportsDistinct"`
	SupportsPagination           bool `json:"supportsPagination"`
	SupportsTrueCurve            bool `json:"supportsTrueCurve"`
	SupportsReturningQueryExtent bool `json:"supportsReturningQueryExtent"`
	SupportsQueryWithDistance    bool `json:"supportsQueryWithDistance"`
	SupportsSqlExpression        bool `json:"supportsSqlExpression"`
}

type LayerConfigDateFieldsTimeReference struct {
	TimeZone               string `json:"timeZone"`
	RespectsDaylightSaving bool   `json:"respectsDaylightSaving"`
}

type LayerConfig struct {
	CurrentVersion                         float32                                           `json:"currentVersion"`
	ID                                     int                                               `json:"id"`
	Name                                   string                                            `json:"name"`
	Description                            string                                            `json:"description"`
	GeometryType                           string                                            `json:"esriGeometryMultipoint"`
	SourceSpatialReference                 LayerConfigSpatialReference                       `json:"sourceSpatialReference"`
	CopyrightText                          string                                            `json:"copyrightText"`
	ParentLayer                            *any                                              `json:"parentLayer"`
	SubLayers                              []any                                             `json:"subLayers"`
	MinScale                               float64                                           `json:"minScale"`
	MaxScale                               float64                                           `json:"maxScale"`
	DrawingInfo                            LayerConfigDrawingInfo                            `json:"drawingInfo"`
	DefaultVisibility                      bool                                              `json:"defaultVisibility"`
	Extent                                 LayerConfigExtent                                 `json:"extent"`
	HasAttachments                         bool                                              `json:"hasAttachments"`
	HTMLPopupType                          string                                            `json:"htmlPopupType"`
	DisplayField                           string                                            `json:"displayField"`
	TypeIDField                            *any                                              `json:"typeIdField"`
	SubtypeFieldName                       *any                                              `json:"subtypeFieldName"`
	SubtypeField                           *any                                              `json:"subtypeField"`
	DefaultSubtypeCode                     *any                                              `json:"defaultSubtypeCode"`
	Fields                                 []LayerConfigField                                `json:"fields"`
	GeometryField                          LayerConfigGeometryField                          `json:"geometryField"`
	Indexes                                []LayerConfigIndex                                `json:"indexes"`
	SubTypes                               []any                                             `json:"subtypes"`
	Relationships                          []any                                             `json:"relationships"`
	CanModifyLayer                         bool                                              `json:"canModifyLayer"`
	CanScaleSymbols                        bool                                              `json:"canScaleSymbols"`
	HasLabels                              bool                                              `json:"hasLabels"`
	Capabilities                           string                                            `json:"capabilities"`
	MaxRecordCount                         int                                               `json:"maxRecordCount"`
	SupportsStatistics                     bool                                              `json:"supportsStatistics"`
	SupportsAdvancedQueries                bool                                              `json:"supportsAdvancedQueries"`
	HasZ                                   bool                                              `json:"hasZ"`
	SupportedQueryFormats                  string                                            `json:"supportedQueryFormats"`
	IsDataVersioned                        bool                                              `json:"isDataVersioned"`
	OwnershipBasedAccessControlForFeatures LayerConfigOwnershipBasedAccessControlForFeatures `json:"ownershipBasedAccessControlForFeatures"`
	UseStandardizedQueries                 bool                                              `json:"useStandardizedQueries"`
	AdvancedQueryCapabilities              LayerConfigAdvancedQueryCapabilities              `json:"advancedQueryCapabilities"`
	SupportsDatumTransformation            bool                                              `json:"supportsDatumTransformation"`
	DateFieldsTimeReference                LayerConfigDateFieldsTimeReference                `json:"dateFieldsTimeReference"`
	SupportsCoordinatesQuantization        bool                                              `json:"supportsCoordinatesQuantization"`
}

func FetchLayerConfig(layerURL *url.URL) (*LayerConfig, error) {
	rawURL := fmt.Sprintf("%s?f=json", layerURL.String())

	page, err := fetch(rawURL)
	if err != nil {
		return nil, err
	}

	var layerConfig LayerConfig
	err = json.Unmarshal(*page, &layerConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal url: %s\n\n%s", rawURL, err))
	}

	return &layerConfig, nil
}
