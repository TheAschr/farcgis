package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"arcserver-go/arcgis"
	"arcserver-go/web"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func createServiceRoute(router *chi.Mux, templates *web.Templates, serviceInfo *arcgis.ServiceInfo) error {
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

func createFolderInfoRoutes(router *chi.Mux, templates *web.Templates, folderInfo *arcgis.FolderInfo) error {
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
		err := createServiceRoute(router, templates, &serviceInfo)
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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	serverURL := expectEnv("SOURCE_ARCGIS_SERVER_URL")
	if serverURL[len(serverURL)-1] == '/' {
		serverURL = serverURL[:len(serverURL)-1]
	}

	router := chi.NewRouter()

	router.Handle("/arcgis/rest/static/*",
		http.StripPrefix("/arcgis/rest/static/",
			http.FileServer(http.Dir("./templates/static")),
		),
	)

	templates := web.Create()

	rootFolderURL, err := url.Parse(fmt.Sprintf("%s/arcgis/rest/services", serverURL))
	if err != nil {
		log.Fatal(err)
	}

	serverInfo, err := arcgis.FetchServerInfo(rootFolderURL)
	if err != nil {
		log.Fatal(err)
	}

	err = createFolderInfoRoutes(router, templates, serverInfo)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Ready...")

	http.ListenAndServe(":8086", router)
}
