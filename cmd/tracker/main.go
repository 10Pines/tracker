package main

import (
	"log"
	"os"

	"github.com/10Pines/tracker/v2/pkg/tracker"
)

func main() {
	apiKey := mustGetAPIKey()
	taskName := mustParseTaskName()
	t := tracker.New(apiKey)
	err := t.CreateBackup(taskName)
	if err != nil {
		log.Fatal(err)
	}
}

func mustGetAPIKey() string {
	key, ok := os.LookupEnv("API_KEY")
	if !ok {
		log.Fatal("API_KEY is missing")
	}
	return key
}

func mustParseTaskName() string {
	if len(os.Args) != 3 || os.Args[1] != "track" {
		log.Fatal("expected format is 'track $TASK_NAME'")
	}
	taskName := os.Args[2]
	if taskName == "" {
		log.Fatal("task name is missing")
	}
	return taskName
}
