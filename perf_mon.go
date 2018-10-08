package main

import (
	"fmt"
	"time"
)

type perfMon struct {
	name    string
	results map[string]time.Duration
}

func (pm *perfMon) track(start time.Time, name string) {
	elapsed := time.Since(start)
	pm.results[name] = elapsed
}

func createMonitoringResultCheckSet(monitors []*perfMon) *checkSet {
	cs := createCheckSet("checks", "emon")
	check := cs.createCheck("slow_checks")
	slow := make(map[string]int)

	for _, pm := range monitors {
		for k, ns := range pm.results {
			key := fmt.Sprintf("%s:%s", pm.name, k)
			if ns > config.EmonSlowCheckThreshold {
				slow[key] = int(ns / time.Millisecond)
			}
		}
	}

	check.Data = slow
	if len(slow) > 0 {
		check.warn(fmt.Sprintf("One or more checks are performing badly (over %v)", config.EmonSlowCheckThreshold))
	}

	return cs
}
