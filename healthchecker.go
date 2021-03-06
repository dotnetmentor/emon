package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/apex/log"
)

func runHealthchecks() ([]*checkSet, int) {
	start := time.Now()

	nodes := getNodes(config.ClusterHTTPEndpoint)
	log.Infof("Running healthchecks on %s", cleanURL(config.ClusterHTTPEndpoint))

	resultChan := make(chan *nodeResult, len(nodes))
	results := make([]*nodeResult, 0)
	checkSets := make([]*checkSet, 0)
	monitors := make([]*perfMon, 0)

	for _, node := range nodes {
		n := node
		go runNodeHealthchecks(n, resultChan)
	}

	c := 0
	for nr := range resultChan {
		c++

		log.Infof("Received result from %s. (%d of %d in %dms)", nr.host, c, len(nodes), int(time.Since(start)/time.Millisecond))
		results = append(results, nr)

		for _, cs := range nr.checkSets {
			checkSets = append(checkSets, cs)
			monitors = append(monitors, cs.monitor)
		}

		if c == len(nodes) {
			close(resultChan)
		}
	}

	cluster := createCheckSet("gossip", "cluster")
	checkSets = append(checkSets, cluster)
	cluster.doClusterConsensusChecks(results)

	checkSets = append(checkSets, createMonitoringResultCheckSet(monitors))

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

func runNodeHealthchecks(node string, resultChan chan *nodeResult) {
	nodeURL, _ := url.Parse(node)
	nodeName := fmt.Sprintf("node-%s", strings.Replace(nodeURL.Hostname(), ".", "-", -1))

	result := &nodeResult{
		host:      nodeURL.Hostname(),
		checkSets: make([]*checkSet, 0),
	}

	gossip := result.createCheckSet("gossip", nodeName)
	stats := result.createCheckSet("stats", nodeName)

	// Do gossip checks
	client := newClient(node)
	gr, err := client.getGossip(gossip)
	if err == nil {
		result.gossip = gr

		gossip.doMasterCount(gr)
		gossip.doSlaveCount(gr)
		gossip.doAliveCount(gr)
	}

	// Do stats checks
	sr, err := client.getStats(stats)
	if err == nil {
		result.stats = sr

		stats.doSysCPUCheck(sr)
		stats.doSysMemoryCheck(sr)
		stats.doProcCPUCheck(sr)
		stats.doProcMemoryCheck(sr)
	}

	resultChan <- result
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
					credentials := ""
					if url.User != nil {
						credentials = fmt.Sprintf("%s@", url.User.String())
					}
					node := fmt.Sprintf("%s://%s%s:%s", url.Scheme, credentials, r, url.Port())
					nodes = append(nodes, node)
				}
			}
		}
	}

	return nodes
}
