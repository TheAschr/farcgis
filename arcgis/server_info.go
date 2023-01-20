package arcgis

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type FolderInfo struct {
	URL      url.URL
	Name     string
	Folders  []FolderInfo
	Services []ServiceInfo

	FolderConfig FolderConfig
}

type ServiceInfo struct {
	URL  url.URL
	Name string
	Type string

	MapServiceConfig      *MapServiceConfig
	FeatureServiceConfig  *FeatureServiceConfig
	GeometryServiceConfig *GeometryServiceConfig
	GPServiceConfig       *GPServiceConfig
}

func removeFolderFromServiceName(serviceName string) string {
	indexOfLastSlash := strings.LastIndex(serviceName, "/")
	if indexOfLastSlash == -1 {
		return serviceName
	}
	return serviceName[indexOfLastSlash+1:]
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
		}
		err = resolveFolderInfo(&subFolderInfo)
		if err != nil {
			return err
		}

		folderInfo.Folders = append(subFolderInfo.Folders, subFolderInfo)
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
			folderInfo.Services = append(folderInfo.Services, ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				MapServiceConfig:      mapServiceConfig,
				FeatureServiceConfig:  nil,
				GeometryServiceConfig: nil,
				GPServiceConfig:       nil,
			})
			break
		case "FeatureServer":
			featureServiceConfig, err := FetchFeatureServiceConfig(serviceURL)
			if err != nil {
				return err
			}
			folderInfo.Services = append(folderInfo.Services, ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				MapServiceConfig:      nil,
				FeatureServiceConfig:  featureServiceConfig,
				GeometryServiceConfig: nil,
				GPServiceConfig:       nil,
			})
			break
		case "GeometryServer":
			geometryServiceConfig, err := FetchGeometryServiceConfig(serviceURL)
			if err != nil {
				return err
			}
			folderInfo.Services = append(folderInfo.Services, ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				MapServiceConfig:      nil,
				FeatureServiceConfig:  nil,
				GeometryServiceConfig: geometryServiceConfig,
				GPServiceConfig:       nil,
			})
			break
		case "GPServer":
			gpServiceConfig, err := FetchGPServiceConfig(serviceURL)
			if err != nil {
				return err
			}
			folderInfo.Services = append(folderInfo.Services, ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				MapServiceConfig:      nil,
				FeatureServiceConfig:  nil,
				GeometryServiceConfig: nil,
				GPServiceConfig:       gpServiceConfig,
			})
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
