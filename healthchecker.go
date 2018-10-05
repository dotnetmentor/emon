package main

import (
	"fmt"
	"time"

	"github.com/apex/log"
)

var monitor perfMon

func runHealthchecks() ([]*checkSet, int) {
	client := newClient(config.ClusterHTTPEndpoint)
	checkSets := make([]*checkSet, 0)
	monitor = perfMon{
		name:    "checks",
		results: make(map[string]time.Duration),
	}

	gossip := createCheckSet("gossip")
	checkSets = append(checkSets, gossip)

	stats := createCheckSet("stats")
	checkSets = append(checkSets, stats)

	// Do gossip checks
	gr, err := client.getGossip(gossip)
	if err == nil {
		gossip.doMasterCount(gr)
		gossip.doSlaveCount(gr)
		gossip.doAliveCount(gr)

		// Do stats checks
		sr, err := client.getStats(stats, gr.ServerIP)
		if err == nil {
			stats.doSysCPUCheck(sr)
			stats.doSysMemoryCheck(sr)
			stats.doProcCPUCheck(sr)
			stats.doProcMemoryCheck(sr)
		}
	}

	checkSets = append(checkSets, monitor.getCheckSet())

	// Output checks
	exitCode := 0

	for _, cs := range checkSets {
		for _, c := range cs.checks {
			topic := "gossip"
			lm := log.WithFields(log.Fields{
				"check":  c.Name,
				"status": c.Status,
				"data":   c.Data,
				"output": c.Output,
			})

			switch c.Status {
			case statusSuccess:
				lm.Info(topic)
				break
			case statusWarning:
				lm.Warn(topic)
				if exitCode == 0 {
					exitCode = 1
				}
				break
			case statusFailed:
				lm.Error(topic)
				exitCode = 2
				break
			}
		}
	}

	return checkSets, exitCode
}

type perfMon struct {
	name    string
	results map[string]time.Duration
}

func (pm *perfMon) track(start time.Time, name string) {
	elapsed := time.Since(start)
	pm.results[name] = elapsed
}

func (pm *perfMon) getCheckSet() *checkSet {
	cs := createCheckSet(pm.name)
	check := cs.createCheck("slow_checks")
	slow := make(map[string]int)

	for k, ns := range pm.results {
		if ns > config.EmonSlowCheckThreshold {
			slow[k] = int(ns / time.Millisecond)
		}
	}

	check.Data = slow
	if len(slow) > 0 {
		check.warn(fmt.Sprintf("One or more checks are performing badly (over %v) %v.", config.EmonSlowCheckThreshold, slow))
	}

	return cs
}
