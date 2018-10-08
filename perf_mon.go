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

func (pm *perfMon) trackCheck(start time.Time, check *check) {
	key := fmt.Sprintf("%s:%s", check.Source, check.Name)
	pm.track(start, key)
}

func (pm *perfMon) getCheckSet() *checkSet {
	cs := createCheckSet(pm.name, "emon")
	check := cs.createCheck("slow_checks")
	slow := make(map[string]int)

	for k, ns := range pm.results {
		if ns > config.EmonSlowCheckThreshold {
			slow[k] = int(ns / time.Millisecond)
		}
	}

	check.Data = slow
	if len(slow) > 0 {
		check.warn(fmt.Sprintf("One or more checks are performing badly (over %v)", config.EmonSlowCheckThreshold))
	}

	return cs
}
