package main

import (
	"log"

	"github.com/10Pines/tracker/v2/internal/http"
	"github.com/10Pines/tracker/v2/internal/logic"
	"github.com/10Pines/tracker/v2/internal/reporter"
	"github.com/10Pines/tracker/v2/internal/schedule"
)

func main() {
	config := mustParseConfig()
	db := mustNewSQL(config.dbPath)

	go schedule.PeriodicallyRunReport(db, reporter.ConsoleReporter)

	l := logic.New(db)
	router := http.NewRouter(l, config.apiKey)

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
