package taskmanager_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager/helpers/errs"
)

type Client struct {
	APIBaseURL string
	HTTPClient *http.Client
	Token      string
	Version    string
	SyncMap    *sync.Map
}

func NewClient(apiBaseURL, token, version string) (*Client, error) {
	return &Client{
		HTTPClient: &http.Client{Timeout: 0},
		APIBaseURL: apiBaseURL,
		Token:      token,
		Version:    version,
		SyncMap:    &sync.Map{},
	}, nil
}

// DoWithLock - Perform API call with lock
func (c *Client) DoWithLock(req *http.Request, key string) ([]byte, error) {
	c.lock(key)
	defer c.unlock(key)
	return c.Do(req)
}

// Do - Perform API call
func (c *Client) Do(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return []byte(""), errs.ErrNoContent
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return body, errs.ErrNotFound
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted {
		var httpErrorResponse HTTPErrorResponse
		err = json.Unmarshal(body, &httpErrorResponse)
		if err == nil && httpErrorResponse.Message != "" {
			return nil, fmt.Errorf("%s: %s", httpErrorResponse.ErrorCode, httpErrorResponse.Message)
		}
		return nil, fmt.Errorf("an error occurred while processing the request\nrequest url: %s\nrequest method: %s\nresponse status: %d\nresponse body: %s", req.URL, req.Method, res.StatusCode, body)
	}

	return body, err
}

// Lock to lock based on key
func (c *Client) lock(key interface{}) {
	mutex := &sync.Mutex{}
	actual, _ := c.SyncMap.LoadOrStore(key, mutex)
	actualMutex := actual.(*sync.Mutex)
	actualMutex.Lock()
	if actualMutex != mutex {
		actualMutex.Unlock()
		c.lock(key)
		return
	}
}

// Unlock to unlock based on key
func (c *Client) unlock(key interface{}) {
	actual, exist := c.SyncMap.Load(key)
	if !exist {
		return
	}
	actualMutex := actual.(*sync.Mutex)
	c.SyncMap.Delete(key)
	actualMutex.Unlock()
}

func (c *Client) Post(url string, body interface{}, lock string) ([]byte, error) {
	bodyContent, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyContent)

	finalUrl := fmt.Sprintf("%s%s", c.APIBaseURL, url)

	req, err := http.NewRequest("POST", finalUrl, bodyReader)
	if err != nil {
		return nil, err
	}

	respBody, err := c.DoWithLock(req, lock)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *Client) Get(url string, lock string) ([]byte, error) {
	finalUrl := fmt.Sprintf("%s%s", c.APIBaseURL, url)
	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.DoWithLock(req, lock)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *Client) Put(url string, body interface{}, lock string) ([]byte, error) {
	bodyContent, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyContent)

	finalUrl := fmt.Sprintf("%s%s", c.APIBaseURL, url)
	req, err := http.NewRequest("PUT", finalUrl, bodyReader)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoWithLock(req, lock)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Patch(url string, body interface{}, lock string) ([]byte, error) {
	bodyContent, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(bodyContent)

	finalUrl := fmt.Sprintf("%s%s", c.APIBaseURL, url)
	req, err := http.NewRequest("PATCH", finalUrl, bodyReader)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoWithLock(req, lock)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Delete(url string, lock string) error {
	finalUrl := fmt.Sprintf("%s%s", c.APIBaseURL, url)
	req, err := http.NewRequest("DELETE", finalUrl, nil)
	if err != nil {
		return err
	}

	_, err = c.DoWithLock(req, lock)
	if err != nil {
		return err
	}

	return nil
}
