package main

import (
	"fmt"
	"math"
	"time"
)

func (client *esHTTPClient) getGossip(set *checkSet) (*gossipResponse, error) {
	defer monitor.track(time.Now(), "collect_gossip")
	gossipCheck := set.createCheck("collect_gossip")

	body, err := client.get("/gossip")
	if err != nil {
		gossipCheck.fail(fmt.Sprintf("An error occured fetching gossip. %s", err))
		return nil, err
	}

	r, err := toGossipResponse(body)
	if err != nil {
		gossipCheck.fail(fmt.Sprintf("An error occured parsing gossip. %s", err))
		return nil, err
	}

	return r, nil
}

func (cs *checkSet) doMasterCount(r *gossipResponse) {
	defer monitor.track(time.Now(), "alive_master")
	check := cs.createCheck("alive_master")

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
	defer monitor.track(time.Now(), "alive_slaves")
	check := cs.createCheck("alive_slaves")

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
	defer monitor.track(time.Now(), "alive_nodes")
	check := cs.createCheck("alive_nodes")

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
