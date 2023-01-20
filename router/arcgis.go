package router

import (
	"encoding/json"
	"errors"
	"farcgis/arcgis"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func createServiceRoutes(
	router *chi.Mux,
	templates *map[string]template.Template,
	serviceInfo *arcgis.ServiceInfo,
) error {
	var err error
	var template template.Template
	var jsonConfig []byte
	var prettyJsonConfig []byte

	switch serviceInfo.Type {
	case "FeatureServer":
		template = (*templates)["feature_service"]

		if jsonConfig, err = json.Marshal(&serviceInfo.FeatureServiceConfig); err != nil {
			return err
		}

		if prettyJsonConfig, err = json.MarshalIndent(&serviceInfo.FeatureServiceConfig, "", "  "); err != nil {
			return err
		}

		break
	case "MapServer":
		template = (*templates)["map_service"]

		if jsonConfig, err = json.Marshal(&serviceInfo.MapServiceConfig); err != nil {
			return err
		}

		if prettyJsonConfig, err = json.MarshalIndent(&serviceInfo.MapServiceConfig, "", "  "); err != nil {
			return err
		}

		break
	case "GeometryServer":
		template = (*templates)["geometry_service"]

		if jsonConfig, err = json.Marshal(&serviceInfo.GeometryServiceConfig); err != nil {
			return err
		}

		if prettyJsonConfig, err = json.MarshalIndent(&serviceInfo.GeometryServiceConfig, "", "  "); err != nil {
			return err
		}

		break
	case "GPServer":
		template = (*templates)["gp_service"]

		if jsonConfig, err = json.Marshal(&serviceInfo.GPServiceConfig); err != nil {
			return err
		}

		if prettyJsonConfig, err = json.MarshalIndent(&serviceInfo.GPServiceConfig, "", "  "); err != nil {
			return err
		}

		break
	default:
		return errors.New(fmt.Sprintf("Unhandled serviceInfo.type '%s'", serviceInfo.Type))
	}

	for _, route := range []string{serviceInfo.URL.Path, fmt.Sprintf("%s/", serviceInfo.URL.Path)} {
		router.Get(route, func(w http.ResponseWriter, r *http.Request) {
			format := r.URL.Query().Get("f")
			switch format {
			case "html":
			case "":
				template.ExecuteTemplate(w, "base.html", &serviceInfo)
				break
			case "json":
				w.Write(jsonConfig)
				break
			case "pjson":
				w.Write(prettyJsonConfig)
				break
			}
		})
	}

	return nil
}

func createFolderInfoRoutes(
	router *chi.Mux,
	templates *map[string]template.Template,
	folderInfo *arcgis.FolderInfo,
) error {
	template := (*templates)["folder"]

	jsonConfig, err := json.Marshal(&folderInfo.FolderConfig)
	if err != nil {
		return err
	}

	prettyJsonConfig, err := json.MarshalIndent(&folderInfo.FolderConfig, "", "  ")
	if err != nil {
		return err
	}

	for _, route := range []string{folderInfo.URL.Path, fmt.Sprintf("%s/", folderInfo.URL.Path)} {
		router.Get(route, func(w http.ResponseWriter, r *http.Request) {
			format := r.URL.Query().Get("f")
			switch format {
			case "html":
			case "":
				template.ExecuteTemplate(w, "base.html", &folderInfo)
				break
			case "json":
				w.Write(jsonConfig)
				break
			case "pjson":
				w.Write(prettyJsonConfig)
				break
			}
		})
	}

	for _, serviceInfo := range folderInfo.Services {
		err := createServiceRoutes(router, templates, &serviceInfo)
		if err != nil {
			return err
		}
	}

	for _, subFolderInfo := range folderInfo.Folders {
		err := createFolderInfoRoutes(router, templates, &subFolderInfo)
		if err != nil {
			return err
		}
	}

	return nil
}
