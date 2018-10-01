package main

import (
	"fmt"
	"math"
)

func (client *esHTTPClient) getGossip(set *checkSet) (*gossipResponse, error) {
	gossipCheck := set.createCheck("server_ip_port")

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

	gossipCheck.Data = fmt.Sprintf("%s:%d", r.ServerIP, r.ServerPort)
	if r.ServerIP == "" || r.ServerPort < 1 {
		gossipCheck.fail("No server ip/port in gossip")
	}

	return r, nil
}

func (cs *checkSet) doMasterCount(r *gossipResponse) {
	count := 0
	for _, m := range r.Members {
		if m.State == "Master" {
			count++
		}
	}

	check := cs.createCheck("exactly_1_master")
	check.Data = fmt.Sprintf("%d master node(s)", count)
	if count != 1 {
		check.fail(fmt.Sprintf("Expected 1 master. Found %d.", count))
	}
}

func (cs *checkSet) doSlaveCount(r *gossipResponse) {
	count := 0
	expected := (config.ClusterSize - 1)
	for _, m := range r.Members {
		if m.State == "Slave" {
			count++
		}
	}

	check := cs.createCheck("exactly_2_slaves")
	check.Data = fmt.Sprintf("%d slave node(s)", count)
	if count != expected {
		check.fail(fmt.Sprintf("Expected %d slaves. Found %d.", expected, count))
	}
}

func (cs *checkSet) doAliveCount(r *gossipResponse) {
	count := 0
	expected := int(math.Ceil(float64(config.ClusterSize) / 2))
	for _, m := range r.Members {
		if m.IsAlive {
			count++
		}
	}

	check := cs.createCheck("alive_nodes")
	check.Data = fmt.Sprintf("%d alive node(s)", count)
	if count < expected {
		check.fail(fmt.Sprintf("Expected at least %d alive nodes. Found %d.", expected, count))
	} else if count < config.ClusterSize {
		check.warn(fmt.Sprintf("Expected %d alive nodes. Found %d.", config.ClusterSize, count))
	}
}
