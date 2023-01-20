package router

import (
	"encoding/json"
	"net/http"

	"farcgis/arcgis"

	"github.com/go-chi/chi/v5"
)

func New(serverInfo *arcgis.FolderInfo) (*chi.Mux, error) {
	router := chi.NewRouter()

	router.Handle("/arcgis/rest/static/*",
		http.StripPrefix("/arcgis/rest/static/",
			http.FileServer(http.Dir("./templates/static")),
		),
	)

	prettyJsonConfig, err := json.MarshalIndent(serverInfo, "", "  ")
	if err != nil {
		return nil, err
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(prettyJsonConfig)
	})

	err = createFolderInfoRoutes(router, createTemplates(), serverInfo)
	if err != nil {
		return nil, err
	}

	return router, err
}
