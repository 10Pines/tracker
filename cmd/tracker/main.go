package main

import (
	"github.com/10Pines/tracker/pkg/tracker"
	"log"
	"os"
	"strconv"
)

func main() {
	key := mustGetApiKey()
	taskID := mustParseTaskID()
	t := tracker.New(key)
	err := t.TrackJob(taskID)
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

func mustParseTaskID() uint {
	if len(os.Args) != 3 || os.Args[1] != "track" {
		log.Fatal("expected format is 'track $TASK_ID'")
	}
	stringID := os.Args[2]
	taskID, err := strconv.ParseInt(stringID, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return uint(taskID)
}
