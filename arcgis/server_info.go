package arcgis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

type FolderInfo struct {
	URL          url.URL       `json:"url"`
	Name         string        `json:"name"`
	Folders      []FolderInfo  `json:"folders"`
	Services     []ServiceInfo `json:"services"`
	ParentFolder *FolderInfo   `json:"-"`

	FolderConfig FolderConfig `json:"folderConfig"`
}

type ServiceInfo struct {
	URL          url.URL     `json:"url"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	ParentFolder *FolderInfo `json:"-"`

	MapServiceConfig      *MapServiceConfig      `json:"mapServiceConfig"`
	FeatureServiceConfig  *FeatureServiceConfig  `json:"featureServiceConfig"`
	GeometryServiceConfig *GeometryServiceConfig `json:"geometryServiceConfig"`
	GPServiceConfig       *GPServiceConfig       `json:"gpServiceConfig"`
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
			folderInfo.Services = append(folderInfo.Services, ServiceInfo{
				URL:                   *serviceURL,
				Name:                  service.Name,
				Type:                  service.Type,
				ParentFolder:          folderInfo,
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
				ParentFolder:          folderInfo,
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
				ParentFolder:          folderInfo,
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
				ParentFolder:          folderInfo,
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

	return &fullDirectory
}

func (folderInfo *FolderInfo) SaveToFile(filename string) error {

	prettyJsonConfig, err := json.MarshalIndent(&folderInfo, "", "  ")
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to marshall folderInfo into json:\n\n%s", err))
	}

	err = ioutil.WriteFile(filename, prettyJsonConfig, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable write folderInfo to file:%s\n\n%s", filename, err))
	}

	return nil
}

func LoadServerInfoFromFile(filename string) (*FolderInfo, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read folderInfo from file: %s\n\n%s", filename, err))
	}

	folderInfo := &FolderInfo{}

	err = json.Unmarshal(file, folderInfo)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to unmarshall folderInfo from file: %s\n\n%s", filename, err))
	}

	return folderInfo, nil
}
