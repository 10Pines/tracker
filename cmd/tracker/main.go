package main

import (
	"log"
	"os"

	"github.com/10Pines/tracker/pkg/tracker"
)

func main() {
	apiKey := mustGetApiKey()
	taskName := mustParseTaskName()
	t := tracker.New(apiKey, tracker.OptionUri("http://localhost:8080/api"))
	err := t.CreateBackup(taskName)
	if err != nil {
		log.Fatal(err)
	}
}

func mustGetApiKey() string {
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
