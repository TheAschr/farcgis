package arcgis

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type LayerInfo struct {
	URL               url.URL      `json:"url"`
	ID                int          `json:"id"`
	Name              string       `json:"name"`
	ParentLayerID     int          `json:"parentLayerId"`
	DefaultVisibility bool         `json:"defaultVisibility"`
	SubLayerIDs       *any         `json:"subLayerIds"`
	MinScale          float64      `json:"minScale"`
	MaxScale          float64      `json:"maxScale"`
	Type              string       `json:"type"`
	GeometryType      string       `json:"geometryType"`
	ParentService     *ServiceInfo `json:"-"`

	LayerConfig *LayerConfig `json:"layerConfig"`
}

type TableInfo struct {
	URL           url.URL      `json:"url"`
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	ParentService *ServiceInfo `json:"-"`

	TableConfig *TableConfig `json:"featureTableConfig"`
}

type ServiceInfo struct {
	URL          url.URL     `json:"url"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	ParentFolder *FolderInfo `json:"-"`
	Layers       []LayerInfo `json:"layers"`
	Tables       []TableInfo `json:"tables"`

	MapServiceConfig      *MapServiceConfig      `json:"mapServiceConfig"`
	FeatureServiceConfig  *FeatureServiceConfig  `json:"featureServiceConfig"`
	GeometryServiceConfig *GeometryServiceConfig `json:"geometryServiceConfig"`
	GPServiceConfig       *GPServiceConfig       `json:"gpServiceConfig"`
}

type FolderInfo struct {
	URL          url.URL       `json:"url"`
	Name         string        `json:"name"`
	Folders      []FolderInfo  `json:"folders"`
	Services     []ServiceInfo `json:"services"`
	ParentFolder *FolderInfo   `json:"-"`

	FolderConfig FolderConfig `json:"folderConfig"`
}

func removeFolderFromServiceName(serviceName string) string {
	indexOfLastSlash := strings.LastIndex(serviceName, "/")
	if indexOfLastSlash == -1 {
		return serviceName
	}
	return serviceName[indexOfLastSlash+1:]
}

func resolveServiceInfo(serviceInfo *ServiceInfo) error {
	// Layers
	var serviceConfigLayers *[]ServiceConfigLayer
	if serviceInfo.FeatureServiceConfig != nil {
		serviceConfigLayers = &serviceInfo.FeatureServiceConfig.Layers
	} else if serviceInfo.MapServiceConfig != nil {
		serviceConfigLayers = &serviceInfo.MapServiceConfig.Layers
	} else if serviceInfo.GeometryServiceConfig != nil {
		serviceConfigLayers = &[]ServiceConfigLayer{}
	} else if serviceInfo.GPServiceConfig != nil {
		serviceConfigLayers = &[]ServiceConfigLayer{}
	} else {
		return errors.New(fmt.Sprintf("Unable to determine serviceConfig type.\nURL: %s", serviceInfo.URL.String()))
	}
	serviceInfo.Layers = make([]LayerInfo, 0)
	for _, layer := range *serviceConfigLayers {
		layerURL, err := url.Parse(fmt.Sprintf("%s/%s", serviceInfo.URL.String(), strconv.Itoa(layer.ID)))
		if err != nil {
			return err
		}
		const expectedLayerType = "Feature Layer"
		if layer.Type != expectedLayerType {
			return errors.New(fmt.Sprintf("Expected layer type to be '%s'. Got '%s'.\nURL: %s", expectedLayerType, layer.Type, serviceInfo.URL.String()))
		}
		layerConfig, err := FetchLayerConfig(layerURL)
		if err != nil {
			return err
		}
		serviceInfo.Layers = append(serviceInfo.Layers, LayerInfo{
			URL:               *layerURL,
			ID:                layer.ID,
			Name:              layer.Name,
			ParentLayerID:     layer.ParentLayerId,
			DefaultVisibility: layer.DefaultVisibility,
			SubLayerIDs:       layer.SubLayerIDs,
			MinScale:          layer.MinScale,
			MaxScale:          layer.MaxScale,
			Type:              layer.Type,
			GeometryType:      layer.GeometryType,
			LayerConfig:       layerConfig,
			ParentService:     serviceInfo,
		})
	}

	// Tables
	var serviceConfigTables *[]ServiceConfigTable
	if serviceInfo.FeatureServiceConfig != nil {
		serviceConfigTables = &serviceInfo.FeatureServiceConfig.Tables
	} else if serviceInfo.MapServiceConfig != nil {
		serviceConfigTables = &serviceInfo.MapServiceConfig.Tables
	} else if serviceInfo.GeometryServiceConfig != nil {
		serviceConfigTables = &[]ServiceConfigTable{}
	} else if serviceInfo.GPServiceConfig != nil {
		serviceConfigTables = &[]ServiceConfigTable{}
	} else {
		return errors.New(fmt.Sprintf("Unable to determine serviceConfig type.\nURL: %s", serviceInfo.URL.String()))
	}
	serviceInfo.Tables = make([]TableInfo, 0)
	for _, table := range *serviceConfigTables {
		tableURL, err := url.Parse(fmt.Sprintf("%s/%s", serviceInfo.URL.String(), strconv.Itoa(table.ID)))
		if err != nil {
			return err
		}

		featureTableConfig, err := FetchTableConfig(tableURL)
		if err != nil {
			return err
		}
		serviceInfo.Tables = append(serviceInfo.Tables, TableInfo{
			URL:           *tableURL,
			ID:            table.ID,
			Name:          table.Name,
			TableConfig:   featureTableConfig,
			ParentService: serviceInfo,
		})
	}

	return nil
}

func resolveFolderInfo(folderInfo *FolderInfo) error {
	// Folders
	folderInfo.Folders = make([]FolderInfo, 0)
	for _, folderName := range folderInfo.FolderConfig.Folders {
		folderUrl, err := url.Parse(fmt.Sprintf("%s/%s", folderInfo.URL.String(), folderName))
		if err != nil {
			return err
		}

		folderConfig, err := FetchFolderConfig(folderUrl)
		if err != nil {
			return err
		}
		subFolderInfo := FolderInfo{
			URL:          *folderUrl,
			Name:         folderName,
			FolderConfig: *folderConfig,
			Folders:      make([]FolderInfo, 0),
			Services:     make([]ServiceInfo, 0),
			ParentFolder: folderInfo,
		}
		err = resolveFolderInfo(&subFolderInfo)
		if err != nil {
			return err
		}

		folderInfo.Folders = append(folderInfo.Folders, subFolderInfo)
	}

	// Services
	folderInfo.Services = make([]ServiceInfo, 0)
	for _, service := range folderInfo.FolderConfig.Services {
		serviceNameWithoutFolder := removeFolderFromServiceName(service.Name)
		serviceURL, err := url.Parse(fmt.Sprintf("%s/%s/%s", folderInfo.URL.String(), serviceNameWithoutFolder, service.Type))
		if err != nil {
			return err
		}
		switch service.Type {
		case "MapServer":
			mapServiceConfig, err := FetchMapServiceConfig(serviceURL)
			if err != nil {
				return err
			}
			serviceInfo := ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				ParentFolder:          folderInfo,
				Layers:                nil,
				Tables:                nil,
				MapServiceConfig:      mapServiceConfig,
				FeatureServiceConfig:  nil,
				GeometryServiceConfig: nil,
				GPServiceConfig:       nil,
			}
			err = resolveServiceInfo(&serviceInfo)
			if err != nil {
				return err
			}
			folderInfo.Services = append(folderInfo.Services, serviceInfo)
			break
		case "FeatureServer":
			featureServiceConfig, err := FetchFeatureServiceConfig(serviceURL)
			if err != nil {
				return err
			}
			serviceInfo := ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				ParentFolder:          folderInfo,
				Layers:                nil,
				Tables:                nil,
				MapServiceConfig:      nil,
				FeatureServiceConfig:  featureServiceConfig,
				GeometryServiceConfig: nil,
				GPServiceConfig:       nil,
			}
			err = resolveServiceInfo(&serviceInfo)
			if err != nil {
				return err
			}
			folderInfo.Services = append(folderInfo.Services, serviceInfo)
			break
		case "GeometryServer":
			geometryServiceConfig, err := FetchGeometryServiceConfig(serviceURL)
			if err != nil {
				return err
			}
			serviceInfo := ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				ParentFolder:          folderInfo,
				Layers:                nil,
				Tables:                nil,
				MapServiceConfig:      nil,
				FeatureServiceConfig:  nil,
				GeometryServiceConfig: geometryServiceConfig,
				GPServiceConfig:       nil,
			}
			err = resolveServiceInfo(&serviceInfo)
			if err != nil {
				return err
			}
			folderInfo.Services = append(folderInfo.Services, serviceInfo)
			break
		case "GPServer":
			gpServiceConfig, err := FetchGPServiceConfig(serviceURL)
			if err != nil {
				return err
			}
			serviceInfo := ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				ParentFolder:          folderInfo,
				Layers:                nil,
				Tables:                nil,
				MapServiceConfig:      nil,
				FeatureServiceConfig:  nil,
				GeometryServiceConfig: nil,
				GPServiceConfig:       gpServiceConfig,
			}
			err = resolveServiceInfo(&serviceInfo)
			if err != nil {
				return err
			}
			folderInfo.Services = append(folderInfo.Services, serviceInfo)
			break
		default:
			return errors.New(fmt.Sprintf("Unhandled service type '%s' at '%s'", service.Type, serviceURL))
		}
	}

	return nil
}

func FetchServerInfo(rootFolderURL *url.URL) (*FolderInfo, error) {
	folderConfig, err := FetchFolderConfig(rootFolderURL)
	if err != nil {
		return nil, err
	}

	folderInfo := &FolderInfo{
		URL:          *rootFolderURL,
		Name:         "services",
		ParentFolder: nil,
		Services:     nil,
		Folders:      nil,
		FolderConfig: *folderConfig,
	}

	err = resolveFolderInfo(folderInfo)
	if err != nil {
		return nil, err
	}

	return folderInfo, nil
}

type Directory struct {
	Label string
	URL   url.URL
}

func (folderInfo *FolderInfo) FullDirectory() *[]Directory {
	fullDirectory := make([]Directory, 0)

	currFolderInfo := folderInfo

	for ok := true; ok; ok = currFolderInfo != nil {
		fullDirectory = append(fullDirectory, Directory{
			Label: currFolderInfo.Name,
			URL:   currFolderInfo.URL,
		})
		currFolderInfo = currFolderInfo.ParentFolder
	}

	ReverseSlice(fullDirectory)

	return &fullDirectory
}

func (serviceInfo *ServiceInfo) FullDirectory() *[]Directory {
	fullDirectory := make([]Directory, 0)

	fullDirectory = append(fullDirectory, Directory{
		Label: serviceInfo.Name,
		URL:   serviceInfo.URL,
	})

	currFolderInfo := serviceInfo.ParentFolder
	if currFolderInfo != nil {
		for ok := true; ok; ok = currFolderInfo != nil {
			fullDirectory = append(fullDirectory, Directory{
				Label: currFolderInfo.Name,
				URL:   currFolderInfo.URL,
			})
			currFolderInfo = currFolderInfo.ParentFolder
		}
	}

	ReverseSlice(fullDirectory)

	return &fullDirectory
}

func (layerInfo *LayerInfo) FullDirectory() *[]Directory {
	fullDirectory := make([]Directory, 0)

	fullDirectory = append(fullDirectory, Directory{
		Label: layerInfo.Name,
		URL:   layerInfo.URL,
	})

	serviceInfo := layerInfo.ParentService

	fullDirectory = append(fullDirectory, Directory{
		Label: serviceInfo.Name,
		URL:   serviceInfo.URL,
	})

	currFolderInfo := serviceInfo.ParentFolder
	if currFolderInfo != nil {
		for ok := true; ok; ok = currFolderInfo != nil {
			fullDirectory = append(fullDirectory, Directory{
				Label: currFolderInfo.Name,
				URL:   currFolderInfo.URL,
			})
			currFolderInfo = currFolderInfo.ParentFolder
		}
	}

	ReverseSlice(fullDirectory)

	return &fullDirectory
}

func (tableInfo *TableInfo) FullDirectory() *[]Directory {
	fullDirectory := make([]Directory, 0)

	fullDirectory = append(fullDirectory, Directory{
		Label: tableInfo.Name,
		URL:   tableInfo.URL,
	})

	serviceInfo := tableInfo.ParentService

	fullDirectory = append(fullDirectory, Directory{
		Label: serviceInfo.Name,
		URL:   serviceInfo.URL,
	})

	currFolderInfo := serviceInfo.ParentFolder
	if currFolderInfo != nil {
		for ok := true; ok; ok = currFolderInfo != nil {
			fullDirectory = append(fullDirectory, Directory{
				Label: currFolderInfo.Name,
				URL:   currFolderInfo.URL,
			})
			currFolderInfo = currFolderInfo.ParentFolder
		}
	}

	ReverseSlice(fullDirectory)

	return &fullDirectory
}

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}
