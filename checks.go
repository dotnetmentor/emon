package main

import (
	"fmt"
	"time"
)

const statusFailed = "failed"
const statusWarning = "warning"
const statusSuccess = "success"

type nodeResult struct {
	host      string
	gossip    *gossipResponse
	checkSets []*checkSet
}

type checkSet struct {
	name    string
	source  string
	checks  []*check
	monitor *perfMon
}

type check struct {
	Name   string      `json:"-"`
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Output string      `json:"output,omitempty"`
}

func (nr *nodeResult) createCheckSet(name string, source string) *checkSet {
	cs := createCheckSet(name, source)
	nr.checkSets = append(nr.checkSets, cs)
	return cs
}

func createCheckSet(name string, source string) *checkSet {
	s := make([]*check, 0)
	m := &perfMon{
		name:    source,
		results: make(map[string]time.Duration),
	}
	return &checkSet{
		name:    name,
		source:  source,
		checks:  s,
		monitor: m,
	}
}

func (s *checkSet) createCheck(name string) *check {
	c := &check{
		Name:   fmt.Sprintf("%s:%s", s.name, name),
		Status: statusSuccess,
	}

	s.checks = append(s.checks, c)

	return c
}

func (s *checkSet) monitorCheck(start time.Time, check *check) {
	s.monitor.track(start, check.Name)
}

func (c *check) fail(reason string) {
	c.Status = statusFailed
	c.Output = reason
}

func (c *check) warn(reason string) {
	c.Status = statusWarning
	c.Output = reason
}
