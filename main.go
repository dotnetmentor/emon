package main

import (
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

func main() {
	log.SetHandler(text.New(os.Stderr))

	client := newClient("http://localhost:12113")
	gossipChecks := createCheckSet("gossip")

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
			"check":  strings.Replace(c.name, "gossip:", "", -1),
			"reason": c.reason,
			"status": c.status,
			"data":   c.data,
		})

		switch c.status {
		case statusSuccess:
			lm.Info(topic)
		case statusWarning:
			lm.Warn(topic)
		case statusFailed:
			lm.Error(topic)
		}

		if c.status != statusSuccess {
			success = false
		}
	}

	exitCode := 0
	if !success {
		exitCode = 1
	}

	os.Exit(exitCode)
}
