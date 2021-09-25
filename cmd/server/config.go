package main

import (
	"log"
	"os"
	"path"
)

type config struct {
	dbPath string
	apiKey string
}

func mustParseConfig() config {
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		log.Fatalf("API_KEY is missing")
	}

	dbPath := os.Getenv("DB_PATH")
	dbPath = path.Join(dbPath, "tracker.db")

	return config{
		dbPath: dbPath,
		apiKey: apiKey,
	}
}
