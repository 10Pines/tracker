package main

import (
	"github.com/10Pines/tracker/internal/http"
	"github.com/10Pines/tracker/internal/logic"
	"github.com/10Pines/tracker/internal/reporter"
	"github.com/10Pines/tracker/internal/schedule"
	"github.com/10Pines/tracker/storage"
	"log"
)

func main() {
	dbPath := storage.DBPath()
	log.Printf("db attached at %s", dbPath)
	db := storage.MustNewSQL(dbPath)

	l := logic.New(db)

	go schedule.PeriodicallyRunReport(db, reporter.ConsoleReporter)

	router := http.NewRouter(l)

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
