package tracker

import (
	"bytes"
	"fmt"
	"net/http"
)

const defaultUri = "https://tracker.10pines.com"

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Tracker struct {
	httpClient httpClient
	uri        string
}

type Option func(*Tracker)

func OptionHttpClient(client httpClient) Option {
	return func(tracker *Tracker) {
		tracker.httpClient = client
	}
}

func OptionUri(uri string) Option {
	return func(tracker *Tracker) {
		tracker.uri = uri
	}
}

func New(options ...Option) *Tracker {
	t := &Tracker{
		httpClient: http.DefaultClient,
		uri:        defaultUri,
	}

	for _, opt := range options {
		opt(t)
	}

	return t
}

func (t Tracker) TrackJob(taskID uint) error {
	var body bytes.Buffer
	url := fmt.Sprintf("%s/tasks/%d/jobs", t.uri, taskID)
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return err
	}
	_, err = t.httpClient.Do(req)
	return err
}
