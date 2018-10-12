package main

type master struct {
	InstanceID     string `json:"instanceId,omitempty"`
	InternalHTTPIP string `json:"internalHttpIp,omitempty"`
	ExternalHTTPIP string `json:"externalHttpIp,omitempty"`
}
