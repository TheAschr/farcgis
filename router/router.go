package router

import (
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

	err := createFolderInfoRoutes(router, createTemplates(), serverInfo)
	if err != nil {
		return nil, err
	}

	return router, err
}
