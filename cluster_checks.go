package main

import (
	"fmt"
	"time"

	"github.com/apex/log"
)

func (cs *checkSet) doClusterConsensusChecks(nodeResults []*nodeResult) {
	cs.doClusterConsensusMasterCheck(nodeResults)
	cs.doClusterConsensusTimeCheck(nodeResults)
}

func (cs *checkSet) doClusterConsensusMasterCheck(nodeResults []*nodeResult) {
	check := cs.createCheck("master_consensus")
	defer cs.monitorCheck(time.Now(), check)

	masters := make([]string, 0)

	for _, nr := range nodeResults {
		if nr.gossip != nil {
			for _, n := range nr.gossip.Members {
				if n.IsAliveMaster() {
					masters = append(masters, n.InstanceID)
				}
			}
		}
	}

	check.Data = distinct(masters)
	if !all(masters, equal) || len(masters) != len(nodeResults) {
		check.fail("Nodes have different masters!")
	}
}

func (cs *checkSet) doClusterConsensusTimeCheck(nodeResults []*nodeResult) {
	check := cs.createCheck("time_consensus")
	defer cs.monitorCheck(time.Now(), check)

	masterTimestamp := time.Now()
	timestamps := make([]time.Time, 0)

	for _, nr := range nodeResults {
		if nr.gossip != nil {
			for _, n := range nr.gossip.Members {
				if n.IsAlive {
					t, err := time.Parse(time.RFC3339, n.Timestamp)
					if err != nil {
						log.Errorf("Failed parsing timestamp from node gossip. Timestamp: %s. Error: %v", n.Timestamp, err)
					}
					timestamps = append(timestamps, t)

					if n.IsAliveMaster() {
						masterTimestamp = t
					}
				}
			}
		}
	}

	check.Data = timestamps
	for _, t := range timestamps {
		diff := t.Sub(masterTimestamp)
		seconds := diff.Seconds()

		maxAhead := 5.0 // TODO: Make configurable
		maxBehind := maxAhead * -1.0
		warnAhead := maxAhead / 2
		warnBehind := warnAhead * -1.0

		if seconds > maxAhead || seconds < maxBehind {
			check.fail(fmt.Sprintf("Clock drift detected between master and one of the nodes (master: %v. diff: %vs.)", masterTimestamp, seconds))
			break
		} else if seconds > warnAhead || seconds < warnBehind {
			check.warn(fmt.Sprintf("Clock is drifting between master and one of the nodes (master: %v. diff: %vs.)", masterTimestamp, seconds))
			break
		}
	}
}

func distinct(a []string) []string {
	results := make([]string, 0)

	for i := 0; i < len(a); i++ {
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
	for i := 0; i < len(a); i++ {
		if equal(a[i], v) {
			return true
		}
	}
	return false
}

func equal(a string, b string) bool {
	return a == b
}
