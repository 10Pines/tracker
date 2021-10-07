package tracker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIKeyHeader is the http header to send the API key
const APIKeyHeader = "X-API-KEY"

const defaultURI = "https://tracker.10pines.com/api"

// CreateBackup defines the necessary arguments required to create a backup
type CreateBackup struct {
	TaskName string `json:"taskName"`
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Tracker is the client used to interact with the system
type Tracker struct {
	httpClient httpClient
	uri        string
	apiKey     string
}

// Option types a generic Tracker customization
type Option func(*Tracker)

// OptionHTTPClient customizes underlying HTTP client
func OptionHTTPClient(client httpClient) Option {
	return func(tracker *Tracker) {
		tracker.httpClient = client
	}
}

// OptionURI customizes the client URI
func OptionURI(uri string) Option {
	return func(tracker *Tracker) {
		tracker.uri = uri
	}
}

// New returns a tracker instance
func New(key string, options ...Option) *Tracker {
	t := &Tracker{
		httpClient: http.DefaultClient,
		uri:        defaultURI,
		apiKey:     key,
	}

	for _, opt := range options {
		opt(t)
	}

	return t
}

// CreateBackup tracks a backup completion
func (t Tracker) CreateBackup(taskName string) error {
	url := fmt.Sprintf("%s/backups", t.uri)
	create := CreateBackup{TaskName: taskName}
	var body bytes.Buffer
	encoder := json.NewEncoder(&body)
	err := encoder.Encode(&create)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, &body)
	req.Header.Set(APIKeyHeader, t.apiKey)
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
