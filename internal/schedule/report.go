package schedule

import (
	"log"
	"time"

	"github.com/10Pines/tracker/v2/internal/logic"
	"github.com/10Pines/tracker/v2/internal/reporter"
)

// PeriodicallyRunReport runs a report at 12PM, every 24 hours
func PeriodicallyRunReport(l logic.Logic, reporter reporter.Reporter) {
	for {
		wait := timeUntilNextReport(time.Now())
		logWaitTime(wait)
		time.Sleep(wait)
		now := time.Now()
		r, err := l.CreateReport(now)
		if err != nil {
			log.Println(err)
			continue
		}
		err = reporter.Process(r)
		if err != nil {
			log.Println(err)
		}
	}
}

func logWaitTime(wait time.Duration) {
	log.Printf("the next report is due in %.2f hs", wait.Round(time.Hour).Hours())
}

func timeUntilNextReport(now time.Time) time.Duration {
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, time.UTC)
	if nextRun.Before(now) {
		nextRun = nextRun.AddDate(0, 0, 1)
	}
	return nextRun.Sub(now)
}
