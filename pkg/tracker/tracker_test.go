package tracker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeClient struct {
	reqs []*http.Request
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	c.reqs = append(c.reqs, req)
	return &http.Response{StatusCode: http.StatusCreated}, nil
}

func newTestTracker(opts ...Option) (*Tracker, *fakeClient) {
	fakeHTTP := &fakeClient{}
	tracker := New("test", append(opts, OptionHTTPClient(fakeHTTP))...)
	return tracker, fakeHTTP
}

func TestTrackerCustomHttpClient(t *testing.T) {
	tracker, fakeHTTP := newTestTracker()

	err := tracker.CreateBackup("payroll weekly backup")
	assert.NoError(t, err)

	assert.Len(t, fakeHTTP.reqs, 1)
}

func TestTrackerDefaultUri(t *testing.T) {
	tracker, fakeHTTP := newTestTracker()

	err := tracker.CreateBackup("payroll weekly backup")
	assert.NoError(t, err)

	assert.Contains(t, fakeHTTP.reqs[0].URL.String(), defaultURI)
}

func TestTrackerCustomUri(t *testing.T) {
	tracker, fakeHTTP := newTestTracker(OptionURI("https://test.com"))

	err := tracker.CreateBackup("payroll weekly backup")
	assert.NoError(t, err)

	assert.Contains(t, fakeHTTP.reqs[0].URL.String(), "https://test.com")
}

func TestTrackerCreateBackupURI(t *testing.T) {
	tracker, fakeHTTP := newTestTracker()

	err := tracker.CreateBackup("payroll weekly backup")
	assert.NoError(t, err)

	assert.Contains(t, fakeHTTP.reqs[0].URL.String(), "/api/backup")
}

func TestTrackerCreateBackupTaskName(t *testing.T) {
	tracker, fakeHTTP := newTestTracker()

	err := tracker.CreateBackup("payroll weekly backup")
	assert.NoError(t, err)

	reqBody, err := ioutil.ReadAll(fakeHTTP.reqs[0].Body)
	assert.NoError(t, err)

	var create CreateBackup
	err = json.Unmarshal(reqBody, &create)
	assert.NoError(t, err)

	assert.Equal(t, create.TaskName, "payroll weekly backup")
}

func TestTrackerCreateBackupApiKey(t *testing.T) {
	tracker, fakeHTTP := newTestTracker()

	err := tracker.CreateBackup("payroll weekly backup")
	assert.NoError(t, err)

	assert.Contains(t, fakeHTTP.reqs[0].Header.Get(APIKeyHeader), "test")
}
