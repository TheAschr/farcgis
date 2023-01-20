package web

import (
	"log"
	"net/url"
)

type Directory struct {
	Label string
	URL   string
}

func expectJoinPath(base string, elem ...string) string {
	joinedPath, err := url.JoinPath(base, elem...)
	if err != nil {
		log.Fatal(err)
	}

	return joinedPath
}

func getDirectoryURLs(directories []*Directory) string {
	result := make([]string, len(directories))
	for i, d := range directories {
		result[i] = d.URL
	}
	if len(result) > 0 {
		return expectJoinPath("/", result[0:]...)
	}
	return "/"
}
