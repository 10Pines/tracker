package main

import (
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/10Pines/tracker/internal/models"
)

type config struct {
	dbPath string
	apiKey string
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
