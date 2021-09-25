package main

import (
	"github.com/10Pines/tracker/internal/http"
	"github.com/10Pines/tracker/internal/logic"
	"github.com/10Pines/tracker/internal/reporter"
	"github.com/10Pines/tracker/internal/schedule"
	"github.com/10Pines/tracker/internal/storage"
	"log"
)

func main() {
	config := mustParseConfig()

	db := storage.MustNewSQL(config.dbPath)

	go schedule.PeriodicallyRunReport(db, reporter.ConsoleReporter)

	l := logic.New(db)
	router := http.NewRouter(l, config.apiKey)

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
