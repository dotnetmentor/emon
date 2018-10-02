package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
)

type esHTTPClient struct {
	BaseURL     string
	Port        int
	Scheme      string
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

type statsResponse struct {
	Proc statsProc `json:"proc"`
	Sys  statsSys  `json:"sys"`
}

type statsProc struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"mem"`
}

type statsSys struct {
	CPU        float64 `json:"cpu"`
	FreeMemory float64 `json:"freeMem"`
}

type errorResponse struct {
}

func newClient(baseURL string) *esHTTPClient {
	u, _ := url.Parse(baseURL)
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		port = 80
	}

	return &esHTTPClient{
		BaseURL:     baseURL,
		Port:        port,
		Scheme:      u.Scheme,
		ContentType: "application/json",
		HTTPClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *esHTTPClient) request(method string, path string, reader io.Reader) ([]byte, error) {
	url := c.BaseURL + path
	if strings.HasPrefix(path, "http") {
		url = path
	}
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

func toStatsResponse(body []byte) (*statsResponse, error) {
	var s = new(statsResponse)
	err := json.Unmarshal(body, &s)
	return s, err
}
