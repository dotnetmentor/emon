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
	nodes := getNodes(config.ClusterHTTPEndpoint)
	log.Infof("Running healthchecks on %s (nodes: %v)", config.ClusterHTTPEndpoint, nodes)

	checkSets := make([]*checkSet, 0)
	monitor = perfMon{
		name:    "checks",
		results: make(map[string]time.Duration),
	}

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

	checkSets = append(checkSets, monitor.getCheckSet())

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

type perfMon struct {
	name    string
	results map[string]time.Duration
}

func (pm *perfMon) track(start time.Time, name string) {
	elapsed := time.Since(start)
	pm.results[name] = elapsed
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
		check.warn(fmt.Sprintf("One or more checks are performing badly (over %v) %v.", config.EmonSlowCheckThreshold, slow))
	}

	return cs
}
