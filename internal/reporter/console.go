package reporter

import (
	"fmt"
	"log"
	"time"

	"github.com/10Pines/tracker/v2/internal/report"
)

func ConsoleReporter(report report.Report) error {
	log.Println("------Status report------")
	log.Printf("Report[timestamp=%s, tasks=%d]", report.Timestamp().Format(time.RFC3339), report.TaskCount())
	for _, taskStatus := range report.Statuses() {
		log.Printf("Task[%s] status=%s", taskStatus.Task.Name, status(taskStatus))
	}
	log.Println("------")
	return nil
}

func status(taskStatus report.TaskStatus) string {
	if !taskStatus.Ready {
		return fmt.Sprintf("INSUFFICIENT_DATAPOINTS[%d/%d]", taskStatus.BackupCount, taskStatus.Task.Datapoints)
	}
	ok := taskStatus.Expected <= taskStatus.BackupCount
	var label string
	if ok {
		label = "OK"
	} else {
		label = "ERR"
	}
	return fmt.Sprintf("%s[%d/%d Tolerance=%d]", label, taskStatus.BackupCount, taskStatus.Task.Datapoints, taskStatus.Task.Tolerance)
}
