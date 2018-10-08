package main

import "time"

func (cs *checkSet) doClusterConsensusChecks(nodeGossip []*gossipResponse) {
	defer monitor.track(time.Now(), "cluster_consensus")
	masters := make([]string, 0)

	for _, ng := range nodeGossip {
		for _, n := range ng.Members {
			if n.IsAliveMaster() {
				masters = append(masters, n.InstanceID)
			}
		}
	}

	masterCheck := cs.createCheck("master_consensus")
	masterCheck.Data = distinct(masters)
	if !all(masters, equal) || len(masters) != len(nodeGossip) {
		masterCheck.fail("Node have different masters!")
	}
}

func distinct(a []string) []string {
	uniques := make([]string, 0)
	uniques = append(uniques, a[0])

	for i := 1; i < len(a); i++ {
		if !equal(a[i], a[0]) {
			uniques = append(uniques, a[i])
		}
	}
	return uniques
}

func all(a []string, f func(a string, b string) bool) bool {
	for i := 1; i < len(a); i++ {
		return f(a[i], a[0])
	}
	return true
}

func equal(a string, b string) bool {
	return a == b
}
