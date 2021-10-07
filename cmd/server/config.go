package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/models"
	"github.com/10Pines/tracker/v2/internal/reporter"
)

type config struct {
	dbPath     string
	apiKey     string
	slackToken string
}

func mustNewSQL(dbPath string) *gorm.DB {
	log.Printf("db attached at %s", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	err = db.AutoMigrate(&models.Backup{}, &models.Task{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func mustParseConfig() config {
	apiKey := mustGetEnv("API_KEY")

	dbPath := os.Getenv("DB_PATH")
	dbPath = path.Join(dbPath, "tracker.db")

	slackToken := mustGetEnv("SLACK_TOKEN")

	return config{
		dbPath:     dbPath,
		apiKey:     apiKey,
		slackToken: slackToken,
	}
}

func mustGetEnv(name string) string {
	apiKey, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf(fmt.Sprintf("%s is missing", name))
	}
	return apiKey
}

func combinedReporter(token string) reporter.Reporter {
	sr := reporter.NewSlackReporter(token, "infraestructura-feed")
	cr := reporter.NewConsoleReporter()
	return reporter.Multiple(cr, sr)
}
