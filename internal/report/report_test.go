package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/10Pines/tracker/v2/internal/models"
)

func TestTaskStatusReady(t *testing.T) {
	tests := []struct {
		name              string
		datapoints        int
		daysOld           int
		expectedReadiness bool
	}{
		{
			name:              "less than datapoints days old",
			datapoints:        5,
			daysOld:           4,
			expectedReadiness: false,
		},
		{
			name:              "same days old as datapoints",
			datapoints:        5,
			daysOld:           5,
			expectedReadiness: true,
		},
		{
			name:              "more than datapoints days old",
			datapoints:        5,
			daysOld:           5,
			expectedReadiness: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := models.NewTask("test", test.datapoints, 0)
			creationTime, _ := time.Parse(time.RFC3339, "2021-02-10T00:00:00Z")
			task.CreatedAt = creationTime
			now := creationTime.AddDate(0, 0, test.daysOld)
			report := newReport(now)
			report.Got(task, 0)
			statuses := report.Statuses()
			assert.Len(t, statuses, 1)
			assert.Equal(t, test.expectedReadiness, statuses[0].Ready)
		})
	}
}
