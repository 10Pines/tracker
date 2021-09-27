package schedule

import (
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/10Pines/tracker/v2/internal/report"
	"github.com/10Pines/tracker/v2/internal/reporter"
)

func PeriodicallyRunReport(db *gorm.DB, reporter reporter.Reporter) {
	//wait := timeUntilNextReport(time.Now())
	//time.Sleep(wait)
	ticker := time.Tick(10 * time.Second)
	for {
		select {
		case <-ticker:
			r, err := report.Run(db)
			if err != nil {
				log.Println(err)
				continue
			}
			err = reporter(r)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func timeUntilNextReport(now time.Time) time.Duration {
	firstRun := time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, time.UTC)
	if firstRun.Before(now) {
		firstRun = firstRun.AddDate(0, 0, 1)
	}
	return firstRun.Sub(now)
}
