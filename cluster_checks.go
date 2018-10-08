package main

import (
	"time"
)

func (cs *checkSet) doClusterConsensusChecks(nodeResults []*nodeResult) {
	check := cs.createCheck("master_consensus")
	defer monitor.trackCheck(time.Now(), check)

	nodeGossip := make([]*gossipResponse, 0)
	for _, nr := range nodeResults {
		nodeGossip = append(nodeGossip, nr.gossip)
	}

	masters := make([]string, 0)

	for _, ng := range nodeGossip {
		for _, n := range ng.Members {
			if n.IsAliveMaster() {
				masters = append(masters, n.InstanceID)
			}
		}
	}

	check.Data = distinct(masters)
	if !all(masters, equal) || len(masters) != len(nodeGossip) {
		check.fail("Nodes have different masters!")
	}
}

func distinct(a []string) []string {
	results := make([]string, 0)
	results = append(results, a[0])

	for i := 1; i < len(a); i++ {
		if !contains(results, a[i]) {
			results = append(results, a[i])
		}
	}

	return results
}

func all(a []string, f func(a string, b string) bool) bool {
	for i := 1; i < len(a); i++ {
		if !f(a[i], a[0]) {
			return false
		}
	}
	return true
}

func contains(a []string, v string) bool {
	for i := 1; i < len(a); i++ {
		if equal(a[i], v) {
			return true
		}
	}
	return false
}

func equal(a string, b string) bool {
	return a == b
}
