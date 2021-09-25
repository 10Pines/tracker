package tracker

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type fakeClient struct {
	reqs []*http.Request
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	c.reqs = append(c.reqs, req)
	return nil, nil
}

func TestTrackerCustomHttpClient(t *testing.T) {
	fakeHttp := &fakeClient{}
	tracker := New(OptionHttpClient(fakeHttp))

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Len(t, fakeHttp.reqs, 1)
}

func TestTrackerDefaultUri(t *testing.T) {
	fakeHttp := &fakeClient{}
	tracker := New(OptionHttpClient(fakeHttp))

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Contains(t, fakeHttp.reqs[0].URL.String(), defaultUri)
}

func TestTrackerCustomUri(t *testing.T) {
	fakeHttp := &fakeClient{}
	tracker := New(OptionHttpClient(fakeHttp), OptionUri("https://test.com"))

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Contains(t, fakeHttp.reqs[0].URL.String(), "https://test.com")
}

func TestTrackerTrackJobURI(t *testing.T) {
	fakeHttp := &fakeClient{}
	tracker := New(OptionHttpClient(fakeHttp))

	err := tracker.TrackJob(1)
	assert.NoError(t, err)

	assert.Contains(t, fakeHttp.reqs[0].URL.String(), "/tasks/1/jobs")
}
