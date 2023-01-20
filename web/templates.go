package web

import (
	"html/template"
)

type Templates = map[string]template.Template

var pages = map[string]string{
	"folder":         "templates/pages/folder.html",
	"featureService": "templates/pages/feature_service.html",
	"mapService":     "templates/pages/map_service.html",
	"layer":          "templates/pages/layer.html",
}

var base = []string{
	"templates/nav_table.html",
	"templates/user_table.html",
	"templates/base.html",
}

var templateFuncs = template.FuncMap{
	"minus": func(a, b int) int {
		return a - b
	},
	"makeFullUrl": func(directories []*Directory, directory *Directory) string {
		fullDirectory := make([]*Directory, 0)
		for _, d := range directories {
			fullDirectory = append(fullDirectory, d)
			if d == directory {
				break
			}
		}
		return getDirectoryURLs(fullDirectory)
	},
}

func Create() *Templates {
	templates := Templates{}
	for pageName, pageFileLoc := range pages {
		files := make([]string, 0)
		files = append(files, pageFileLoc)
		files = append(files, base...)

		templates[pageName] = *template.Must(
			template.New("").Funcs(templateFuncs).ParseFiles(
				files...,
			),
		)
	}

	return &templates
}
