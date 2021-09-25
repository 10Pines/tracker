package schedule

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestTimeTillNextRun(t *testing.T) {
	tests := []struct {
		name             string
		now              time.Time
		expectedDuration time.Duration
	}{
		{
			name:             "before midnight",
			now:              date(10, 10, 10, 0),
			expectedDuration: parseDuration("14h"),
		},
		{
			name:             "at midnight",
			now:              date(10, 10, 24, 0),
			expectedDuration: parseDuration("24h"),
		},
		{
			name:             "after midnight",
			now:              date(10, 11, 10, 0),
			expectedDuration: parseDuration("14h"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wait := timeUntilNextReport(test.now)
			assert.Equal(t, test.expectedDuration, wait)
		})
	}
}

func date(month time.Month, day, hour, min int) time.Time {
	return time.Date(2021, month, day, hour, min, 0, 0, time.UTC)
}

func parseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err)
	}
	return d
}
