package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/models"
	"github.com/10Pines/tracker/v2/internal/reporter"
)

type config struct {
	apiKey     string
	slackToken string
	dbDSN      string
}

func mustNewSQL(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	err = db.AutoMigrate(&models.Backup{}, &models.Task{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func mustParseConfig() config {
	apiKey := mustGetEnv("API_KEY")
	slackToken := mustGetEnv("SLACK_TOKEN")
	dbDSN := mustGetEnv("DB_DSN")

	return config{
		apiKey:     apiKey,
		slackToken: slackToken,
		dbDSN:      dbDSN,
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
