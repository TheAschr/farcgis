package arcgis

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetch(rawURL string) (*[]byte, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to scrape url: %s\n\n%s", rawURL, err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to ready body of url: %s\n\n%s", rawURL, err))
	}

	return &body, nil
}
