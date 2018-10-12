package main

import "encoding/json"

type gossipResponse struct {
	ServerIP   string         `json:"serverIp"`
	ServerPort int            `json:"serverPort"`
	Members    []gossipMember `json:"members"`
}

type gossipMember struct {
	InstanceID     string `json:"instanceId"`
	State          string `json:"state"`
	IsAlive        bool   `json:"isAlive"`
	Timestamp      string `json:"timestamp"`
	InternalHTTPIP string `json:"internalHttpIp"`
	ExternalHTTPIP string `json:"externalHttpIp"`
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

func (m gossipMember) IsAliveMaster() bool {
	return m.IsAlive && m.State == "Master"
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
