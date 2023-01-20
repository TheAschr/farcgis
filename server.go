package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"farcgis/arcgis"
	"farcgis/router"

	"github.com/joho/godotenv"
)

func expectEnv(envName string) string {
	baseUrl := os.Getenv(envName)

	if baseUrl == "" {
		log.Fatalf("Missing expected environment variable '%s'", envName)
	}

	return baseUrl
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env file:%s", err)
	}

	serverURL := expectEnv("SOURCE_ARCGIS_SERVER_URL")
	if serverURL[len(serverURL)-1] == '/' {
		serverURL = serverURL[:len(serverURL)-1]
	}

	portEnv := expectEnv("PORT")
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Fatalf("Invalid port: %s\n\n%s", portEnv, err)
	}

	rootFolderURL, err := url.Parse(fmt.Sprintf("%s/arcgis/rest/services", serverURL))
	if err != nil {
		log.Fatalf("Invalid root folder url:\n\n%s", err)
	}

	const serverInfoFilename = "server-info.json"

	var serverInfo *arcgis.FolderInfo

	_, err = os.Stat(serverInfoFilename)
	if err == nil {
		serverInfo, err = arcgis.LoadServerInfoFromFile(serverInfoFilename)
		if err != nil {
			log.Fatal(err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		serverInfo, err = arcgis.FetchServerInfo(rootFolderURL)
		if err != nil {
			log.Fatal(err)
		}

		err = serverInfo.SaveToFile(serverInfoFilename)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}

	router, err := router.New(serverInfo)
	if err != nil {
		log.Fatalf("Unable to create router:\n\n%s", err)
	}

	addr := fmt.Sprintf(":%s", strconv.Itoa(port))

	log.Printf("Server running on %s\n", addr)

	http.ListenAndServe(addr, router)
}
