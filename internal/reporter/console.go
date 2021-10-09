package reporter

import (
	"fmt"
	"log"
	"time"

	"github.com/10Pines/tracker/v2/internal/logic"
)

type consoleReporter struct {
}

// NewConsoleReporter creates a Reporter instance that uses STDOUT
func NewConsoleReporter() Reporter {
	return consoleReporter{}
}

// Name returns the reporter Name
func (c consoleReporter) Name() string {
	return "console"
}

// Process prints the report results using STDOUT
func (c consoleReporter) Process(report logic.Report) error {
	log.Println("------Status report------")
	log.Printf("Report[timestamp=%s, tasks=%d]", report.Timestamp.Format(time.RFC3339), report.TaskCount())
	for _, taskStatus := range report.Statuses() {
		log.Printf("Task[%s] status=%s", taskStatus.Task.Name, status(taskStatus))
	}
	log.Println("------")
	return nil
}

func status(taskStatus logic.TaskStatus) string {
	if !taskStatus.Ready {
		return fmt.Sprintf("INSUFFICIENT_DATAPOINTS[%d/%d]", taskStatus.BackupCount, taskStatus.Task.Datapoints)
	}
	var label string
	if taskStatus.IsOK() {
		label = "OK"
	} else {
		label = "ERR"
	}
	return fmt.Sprintf("%s[%d/%d Tolerance=%d]", label, taskStatus.BackupCount, taskStatus.Task.Datapoints, taskStatus.Task.Tolerance)
}
