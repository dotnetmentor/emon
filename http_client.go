package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/apex/log"
)

type esHTTPClient struct {
	BaseURL     string
	ContentType string
	HTTPClient  *http.Client
}

type gossipResponse struct {
	ServerIP   string         `json:"serverIp"`
	ServerPort int            `json:"serverPort"`
	Members    []gossipMember `json:"members"`
}

type gossipMember struct {
	InstanceID string `json:"instanceId"`
	State      string `json:"state"`
	IsAlive    bool   `json:"isAlive"`
}

type errorResponse struct {
}

func newClient(baseURL string) *esHTTPClient {
	return &esHTTPClient{
		BaseURL:     baseURL,
		ContentType: "application/json",
		HTTPClient:  &http.Client{},
	}
}

func (c *esHTTPClient) request(method string, path string, reader io.Reader) ([]byte, error) {
	url := c.BaseURL + path
	req, _ := http.NewRequest(method, url, reader)
	req.Header.Set("Content-Type", c.ContentType)

	log.Debugf("[DEBUG] HTTP %s %s", method, url)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s - %s", method, url, resp.Status)
	}

	return body, err
}

func (c *esHTTPClient) get(path string) ([]byte, error) {
	return c.request("GET", path, nil)
}

func (c *esHTTPClient) post(path string, body io.Reader) ([]byte, error) {
	return c.request("POST", path, body)
}

func (c *esHTTPClient) delete(path string) ([]byte, error) {
	return c.request("DELETE", path, nil)
}

func toGossipResponse(body []byte) (*gossipResponse, error) {
	var s = new(gossipResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}
