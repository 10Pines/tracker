package tracker

import (
	"bytes"
	"errors"
	"fmt"
	trackerhttp "github.com/10Pines/tracker/internal/http"
	"io/ioutil"
	"net/http"
)

const defaultUri = "https://tracker.10pines.com/api"

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Tracker struct {
	httpClient httpClient
	uri        string
	apiKey     string
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

func New(key string, options ...Option) *Tracker {
	t := &Tracker{
		httpClient: http.DefaultClient,
		uri:        defaultUri,
		apiKey:     key,
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
	req.Header.Set(trackerhttp.ApiKeyHeader, t.apiKey)
	if err != nil {
		return err
	}
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		errorMsg := fmt.Sprintf("failed with status=%d, body=%s", resp.StatusCode, body)
		return errors.New(errorMsg)
	}
	return nil
}
