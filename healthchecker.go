package main

import (
	"strings"

	"github.com/apex/log"
)

func runHealthchecks() ([]*checkSet, int) {
	client := newClient(config.ClusterHTTPEndpoint)
	checks := make([]*checkSet, 0)

	gossipChecks := createCheckSet("gossip")
	checks = append(checks, gossipChecks)

	// Do gossip checks
	r, err := client.getGossip(gossipChecks)
	if err == nil {
		gossipChecks.doMasterCount(r)
		gossipChecks.doSlaveCount(r)
		gossipChecks.doAliveCount(r)
	}

	// Output checks
	success := true
	for _, c := range gossipChecks.checks {
		topic := "gossip"
		lm := log.WithFields(log.Fields{
			"check":  strings.Replace(c.Name, "gossip:", "", -1),
			"reason": c.Reason,
			"status": c.Status,
			"data":   c.Data,
			"output": c.Output,
		})

		switch c.Status {
		case statusSuccess:
			lm.Info(topic)
		case statusWarning:
			lm.Warn(topic)
		case statusFailed:
			lm.Error(topic)
		}

		if c.Status != statusSuccess {
			success = false
		}
	}

	exitCode := 0
	if !success {
		exitCode = 1
	}

	return checks, exitCode
}
