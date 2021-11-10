package shared

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
			report := NewReport(now)
			report.Got(task, BackupStats{})
			statuses := report.Statuses()
			assert.Len(t, statuses, 1)
			assert.Equal(t, test.expectedReadiness, statuses[0].Ready)
		})
	}
}

func TestReport_IsOk_allTasksOk(t *testing.T) {
	report := NewReport(time.Now())
	task := models.NewTask("test", 5, 2)
	report.Got(task, BackupStats{CountWithinDatapoints: 5})
	assert.True(t, report.IsOK())
}

func TestReport_IsOk_tasksInAlarmState(t *testing.T) {
	report := NewReport(time.Now())
	task := models.NewTask("test", 5, 2)
	report.Got(task, BackupStats{CountWithinDatapoints: 5})
	report.Got(task, BackupStats{CountWithinDatapoints: 0})
	assert.False(t, report.IsOK())
}

func TestTaskStatus_IsOk_NotReady(t *testing.T) {
	ts := TaskStatus{
		Ready: false,
	}
	assert.True(t, ts.IsOK())
}

func TestTaskStatus_IsOk_Ready_UnderThreshold(t *testing.T) {
	ts := TaskStatus{
		Ready:       true,
		BackupCount: 5,
		Expected:    5,
	}
	assert.True(t, ts.IsOK())
}

func TestTaskStatus_IsOk_Ready_OverThreshold(t *testing.T) {
	ts := TaskStatus{
		Ready:       true,
		BackupCount: 2,
		Expected:    5,
	}
	assert.False(t, ts.IsOK())
}
