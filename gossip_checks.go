package main

import (
	"fmt"
	"math"
	"time"
)

func (client *esHTTPClient) getGossip(cs *checkSet) (*gossipResponse, error) {
	check := cs.createCheck("collect_gossip")
	defer cs.monitorCheck(time.Now(), check)

	body, err := client.get("/gossip")
	if err != nil {
		check.fail(fmt.Sprintf("An error occured fetching gossip. %s", err))
		return nil, err
	}

	r, err := toGossipResponse(body)
	if err != nil {
		check.fail(fmt.Sprintf("An error occured parsing gossip. %s", err))
		return nil, err
	}

	return r, nil
}

func (cs *checkSet) doMasterCount(r *gossipResponse) {
	check := cs.createCheck("alive_master")
	defer cs.monitorCheck(time.Now(), check)

	count := 0
	for _, m := range r.Members {
		if m.IsAliveMaster() {
			count++
		}
	}

	check.Data = count
	check.Output = fmt.Sprintf("%d master node(s)", count)
	if count != 1 {
		check.fail(fmt.Sprintf("Expected 1 master. Found %d.", count))
	}
}

func (cs *checkSet) doSlaveCount(r *gossipResponse) {
	check := cs.createCheck("alive_slaves")
	defer cs.monitorCheck(time.Now(), check)

	count := 0
	failLevel := int(math.Ceil(float64(config.ClusterSize)/2)) - 1
	warnLevel := config.ClusterSize - 1
	for _, m := range r.Members {
		if m.State == "Slave" && m.IsAlive {
			count++
		}
	}

	check.Data = count
	check.Output = fmt.Sprintf("%d slave node(s)", count)
	if count < failLevel {
		check.fail(fmt.Sprintf("Expected at least %d slave(s). Found %d.", failLevel, count))
	} else if count < warnLevel {
		check.warn(fmt.Sprintf("Want %d or more slave(s). Found %d.", warnLevel, count))
	}
}

func (cs *checkSet) doAliveCount(r *gossipResponse) {
	check := cs.createCheck("alive_nodes")
	defer cs.monitorCheck(time.Now(), check)

	count := 0
	failLevel := int(math.Ceil(float64(config.ClusterSize) / 2))
	warnLevel := config.ClusterSize
	for _, m := range r.Members {
		if m.IsAlive {
			count++
		}
	}

	check.Data = count
	check.Output = fmt.Sprintf("%d alive node(s)", count)
	if count < failLevel {
		check.fail(fmt.Sprintf("Expected at least %d alive node(s). Found %d.", failLevel, count))
	} else if count < warnLevel {
		check.warn(fmt.Sprintf("Want %d or more alive node(s). Found %d.", warnLevel, count))
	}
}
