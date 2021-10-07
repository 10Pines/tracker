package reporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeFormat(t *testing.T) {
	date := time.Date(2021, 10, 8, 15, 23, 0, 0, time.UTC)
	assert.Equal(t, "2021-10-08 15:23", date.Format(shortISO))
}
