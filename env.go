package main

import (
	"log"
	"os"
)

func expectEnv(envName string) string {
	baseUrl := os.Getenv(envName)

	if baseUrl == "" {
		log.Fatalf("Missing expected environment variable '%s'", envName)
	}

	return baseUrl
}
