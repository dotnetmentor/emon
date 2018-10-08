package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/apex/log"
)

var monitor perfMon

func runHealthchecks() ([]*checkSet, int) {
	start := time.Now()

	nodes := getNodes(config.ClusterHTTPEndpoint)
	log.Infof("Running healthchecks on %s", config.ClusterHTTPEndpoint)

	checkSets := make([]*checkSet, 0)
	monitor = perfMon{
		name:    "checks",
		results: make(map[string]time.Duration),
	}

	nodeGossip := make([]*gossipResponse, 0)

	for _, node := range nodes {
		nodeURL, _ := url.Parse(node)
		nodeName := fmt.Sprintf("node-%s", strings.Replace(nodeURL.Hostname(), ".", "-", -1))
		gossip := createCheckSet("gossip", nodeName)
		checkSets = append(checkSets, gossip)

		stats := createCheckSet("stats", nodeName)
		checkSets = append(checkSets, stats)

		// Do gossip checks
		client := newClient(node)
		gr, err := client.getGossip(gossip)
		if err == nil {
			nodeGossip = append(nodeGossip, gr)

			gossip.doMasterCount(gr)
			gossip.doSlaveCount(gr)
			gossip.doAliveCount(gr)

			// Do stats checks
			sr, err := client.getStats(stats)
			if err == nil {
				stats.doSysCPUCheck(sr)
				stats.doSysMemoryCheck(sr)
				stats.doProcCPUCheck(sr)
				stats.doProcMemoryCheck(sr)
			}
		}
	}

	cluster := createCheckSet("gossip", "cluster")
	checkSets = append(checkSets, cluster)

	cluster.doClusterConsensusChecks(nodeGossip)

	checkSets = append(checkSets, monitor.getCheckSet())

	log.Infof("Completed all checks in %dms", int(time.Since(start)/time.Millisecond))

	// Output checks
	exitCode := 0

	for _, cs := range checkSets {
		for _, c := range cs.checks {
			topic := "gossip"
			lm := log.WithFields(log.Fields{
				"check":  c.Name,
				"source": cs.source,
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

func getNodes(endpoint string) []string {
	nodes := make([]string, 0)

	url, err := url.Parse(endpoint)
	if err != nil {
		log.Errorf("Failed parsing url %v. Error: %v", endpoint, err)
	} else {
		results, err := net.LookupHost(url.Hostname())
		if err != nil {
			log.Errorf("Failed looking up endpoint %v. Error: %v", endpoint, err)
		} else {
			for _, r := range results {
				if r != "::1" {
					node := fmt.Sprintf("%s://%s:%s", url.Scheme, r, url.Port())
					nodes = append(nodes, node)
				}
			}
		}
	}

	return nodes
}
