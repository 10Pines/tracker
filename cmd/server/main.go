package main

import (
	"log"

	"github.com/10Pines/tracker/v2/internal/http"
	"github.com/10Pines/tracker/v2/internal/logic"
)

func main() {
	config := mustParseConfig()
	db := mustNewSQL(config.dbDSN)
	reporter := combinedReporter(config.slackToken)

	l := logic.New(db, reporter)

	router := http.NewRouter(l, config.apiKey)

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
