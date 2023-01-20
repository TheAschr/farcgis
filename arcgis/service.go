package arcgis

type ServiceConfigLayer struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	ParentLayerId     int     `json:"parentLayerId"`
	DefaultVisibility bool    `json:"defaultVisibility"`
	SubLayerIDs       *any    `json:"subLayerIds"`
	MinScale          float64 `json:"minScale"`
	MaxScale          float64 `json:"maxScale"`
	Type              string  `json:"type"`
	GeometryType      string  `json:"geoemtryType"`
}

type ServiceConfigTable struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // This might not exist?
}

type ServiceConfigExtent struct {
	XMin             float64                       `json:"xmin"`
	YMin             float64                       `json:"ymin"`
	XMax             float64                       `json:"xmax"`
	YMax             float64                       `json:"ymax"`
	SpatialReference ServiceConfigSpatialReference `json:"spatialReference"`
}

type ServiceConfigSpatialReference struct {
	WKID       int `json:"wkid"`
	LatestWKID int `json:"latestWkid"`
}

type ServiceConfigDatumTransformationGeoTransform struct {
	WKID             int    `json:"wkid"`
	LatestWKID       int    `json:"latestWkid"`
	TransformForward bool   `json:"transformForward"`
	Name             string `json:"name"`
}

type ServiceConfigDatumTransformation struct {
	GeoTransforms []ServiceConfigDatumTransformationGeoTransform `json:"geoTransforms"`
}
