package main

import "fmt"

const statusFailed = "failed"
const statusWarning = "warning"
const statusSuccess = "success"

type checkSet struct {
	name   string
	checks []*check
}

type check struct {
	name   string
	status string
	reason string
	data   string
}

func createCheckSet(name string) *checkSet {
	s := make([]*check, 0)
	return &checkSet{
		name:   name,
		checks: s,
	}
}

func (s *checkSet) createCheck(name string) *check {
	c := &check{
		name:   fmt.Sprintf("%s:%s", s.name, name),
		status: statusSuccess,
	}

	s.checks = append(s.checks, c)

	return c
}

func (c *check) fail(reason string) {
	c.status = statusFailed
	c.reason = reason
}

func (c *check) warn(reason string) {
	c.status = statusWarning
	c.reason = reason
}
