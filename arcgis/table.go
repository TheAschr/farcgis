package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type TableConfigFieldDomainCodedValue struct {
	Name string      `json:"name"`
	Code interface{} `json:"code"` // Can be string or number
}

type TableConfigFieldDomain struct {
	Type        string                             `json:"type"`
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	CodedValues []TableConfigFieldDomainCodedValue `json:"codedValues"`
	MergePolicy string                             `json:"mergePolicy"`
	SplitPolicy string                             `json:"splitPolicy"`
}

type TableConfigField struct {
	Name   string                  `json:"name"`
	Type   string                  `json:"type"`
	Alias  string                  `json:"alias"`
	Domain *TableConfigFieldDomain `json:"domain"`
}

type TableConfigIndex struct {
	Name        string `json:"name"`
	Fields      string `json:"fields"`
	IsAscending bool   `json:"isAscending"`
	IsUnique    bool   `json:"isUnique"`
	Description string `json:"description"`
}

type TableConfigAdvancedQueryCapabilities struct {
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

type TableConfig struct {
	CurrentVersion            float32                              `json:"currentVersion"`
	ID                        int                                  `json:"id"`
	Name                      string                               `json:"name"`
	Type                      string                               `json:"type"`
	Description               *string                              `json:"description"`
	HasAttachments            bool                                 `json:"hasAttachments"`
	HTMLPopupType             string                               `json:"htmlPopupType"`
	DisplayField              string                               `json:"displayField"`
	TypeIDField               *any                                 `json:"typeIdField"`
	SubtypeFieldName          *any                                 `json:"subtypeFieldName"`
	SubtypeField              *any                                 `json:"subtypeField"`
	DefaultSubtypeCode        *any                                 `json:"defaultSubtypeCode"`
	Fields                    []TableConfigField                   `json:"fields"`
	Indexes                   []TableConfigIndex                   `json:"indexes"`
	Subtypes                  []any                                `json:"subtypes"`
	Relationships             []any                                `json:"relationships"`
	Capabilities              string                               `json:"capabilities"`
	MaxRecordCount            int                                  `json:"maxRecordCount"`
	SupportsStatistics        bool                                 `json:"supportsStatistics"`
	SupportsAdvancedQueries   bool                                 `json:"supportsAdvancedQueries"`
	SupportedQueryFormats     string                               `json:"supportedQueryFormats"`
	IsDataVersioned           bool                                 `json:"isDataVersioned"`
	UseStandardizedQueries    bool                                 `json:"useStandardizedQueries"`
	AdvancedQueryCapabilities TableConfigAdvancedQueryCapabilities `json:"advancedQueryCapabilities"`
}

func FetchTableConfig(tableURL *url.URL) (*TableConfig, error) {
	rawURL := fmt.Sprintf("%s?f=json", tableURL.String())

	page, err := fetch(rawURL)
	if err != nil {
		return nil, err
	}

	var tableConfig TableConfig
	err = json.Unmarshal(*page, &tableConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal url: %s\n\n%s", rawURL, err))
	}

	return &tableConfig, nil
}
