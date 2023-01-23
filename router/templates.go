package router

import (
	"html/template"
)

var pages = map[string]string{
	"folder":         "templates/pages/folder.html",
	"featureService": "templates/pages/feature_service.html",
	"mapService":     "templates/pages/map_service.html",
	"layer":          "templates/pages/layer.html",
	"table":          "templates/pages/table.html",
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
}

func createTemplates() *map[string]template.Template {
	templates := map[string]template.Template{}
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
