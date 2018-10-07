package main

import "fmt"

const statusFailed = "failed"
const statusWarning = "warning"
const statusSuccess = "success"

type checkSet struct {
	name   string
	source string
	checks []*check
}

type check struct {
	Name   string      `json:"-"`
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Output string      `json:"output,omitempty"`
}

func createCheckSet(name string, source string) *checkSet {
	s := make([]*check, 0)
	return &checkSet{
		name:   name,
		source: source,
		checks: s,
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

func (c *check) fail(reason string) {
	c.Status = statusFailed
	c.Output = reason
}

func (c *check) warn(reason string) {
	c.Status = statusWarning
	c.Output = reason
}
