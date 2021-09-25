package tracker

import (
	trackerhttp "github.com/10Pines/tracker/internal/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type fakeClient struct {
	reqs []*http.Request
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	c.reqs = append(c.reqs, req)
	return &http.Response{StatusCode: http.StatusCreated}, nil
}

func newTestTracker(opts ...Option) (*Tracker, *fakeClient) {
	fakeHttp := &fakeClient{}
	tracker := New("test", append(opts, OptionHttpClient(fakeHttp))...)
	return tracker, fakeHttp
}

func TestTrackerCustomHttpClient(t *testing.T) {
	tracker, fakeHttp := newTestTracker()

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Len(t, fakeHttp.reqs, 1)
}

func TestTrackerDefaultUri(t *testing.T) {
	tracker, fakeHttp := newTestTracker()

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Contains(t, fakeHttp.reqs[0].URL.String(), defaultUri)
}

func TestTrackerCustomUri(t *testing.T) {
	tracker, fakeHttp := newTestTracker(OptionUri("https://test.com"))

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Contains(t, fakeHttp.reqs[0].URL.String(), "https://test.com")
}

func TestTrackerTrackJobURI(t *testing.T) {
	tracker, fakeHttp := newTestTracker()

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Contains(t, fakeHttp.reqs[0].URL.String(), "/api/tasks/1/jobs")
}

func TestTrackerTrackJobApiKey(t *testing.T) {
	tracker, fakeHttp := newTestTracker()

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Contains(t, fakeHttp.reqs[0].Header.Get(trackerhttp.ApiKeyHeader), "test")
}
