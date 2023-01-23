package router

import (
	"encoding/json"
	"errors"
	"farcgis/arcgis"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func createLayerRoutes(
	router *chi.Mux,
	templates *map[string]template.Template,
	layerInfo *arcgis.LayerInfo,
) error {
	template := (*templates)["layer"]

	jsonConfig, err := json.Marshal(&layerInfo.LayerConfig)
	if err != nil {
		return err
	}

	prettyJsonConfig, err := json.MarshalIndent(&layerInfo.LayerConfig, "", "  ")
	if err != nil {
		return err
	}

	for _, route := range []string{layerInfo.URL.Path, fmt.Sprintf("%s/", layerInfo.URL.Path)} {
		router.Get(route, func(w http.ResponseWriter, r *http.Request) {
			format := r.URL.Query().Get("f")
			switch format {
			case "html":
			case "":
				err = template.ExecuteTemplate(w, "base.html", layerInfo)
				if err != nil {
					log.Fatal(err)
				}
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

func createTableRoutes(
	router *chi.Mux,
	templates *map[string]template.Template,
	tableInfo *arcgis.TableInfo,
) error {
	template := (*templates)["table"]

	jsonConfig, err := json.Marshal(&tableInfo.TableConfig)
	if err != nil {
		return err
	}

	prettyJsonConfig, err := json.MarshalIndent(&tableInfo.TableConfig, "", "  ")
	if err != nil {
		return err
	}

	for _, route := range []string{tableInfo.URL.Path, fmt.Sprintf("%s/", tableInfo.URL.Path)} {
		router.Get(route, func(w http.ResponseWriter, r *http.Request) {
			format := r.URL.Query().Get("f")
			switch format {
			case "html":
			case "":
				err = template.ExecuteTemplate(w, "base.html", tableInfo)
				if err != nil {
					log.Fatal(err)
				}
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
		template = (*templates)["featureService"]

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
				err = template.ExecuteTemplate(w, "base.html", serviceInfo)
				if err != nil {
					log.Fatal(err)
				}
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

	for _, layerInfo := range serviceInfo.Layers {
		layerInfo := layerInfo
		err = createLayerRoutes(router, templates, &layerInfo)
		if err != nil {
			return err
		}
	}

	for _, tableInfo := range serviceInfo.Tables {
		tableInfo := tableInfo
		err = createTableRoutes(router, templates, &tableInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func createFolderInfoRoutes(
	router *chi.Mux,
	templates *map[string]template.Template,
	folderInfo *arcgis.FolderInfo,
) error {
	template := (*templates)["folder"]

	jsonConfig, err := json.Marshal(folderInfo.FolderConfig)
	if err != nil {
		return err
	}

	prettyJsonConfig, err := json.MarshalIndent(folderInfo.FolderConfig, "", "  ")
	if err != nil {
		return err
	}

	for _, route := range []string{folderInfo.URL.Path, fmt.Sprintf("%s/", folderInfo.URL.Path)} {
		router.Get(route, func(w http.ResponseWriter, r *http.Request) {
			format := r.URL.Query().Get("f")
			switch format {
			case "html":
			case "":
				err = template.ExecuteTemplate(w, "base.html", folderInfo)
				if err != nil {
					log.Fatal(err)
				}
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
		serviceInfo := serviceInfo // Create a copy so that it doesn't change in route handler
		if err := createServiceRoutes(router, templates, &serviceInfo); err != nil {
			return err
		}
	}

	for _, subFolderInfo := range folderInfo.Folders {
		subFolderInfo := subFolderInfo // Create a copy so that it doesn't change in route handler
		if err := createFolderInfoRoutes(router, templates, &subFolderInfo); err != nil {
			return err
		}
	}

	return nil
}
