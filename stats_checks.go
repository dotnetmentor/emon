package main

import (
	"fmt"
	"math"
	"time"
)

func (client *esHTTPClient) getStats(cs *checkSet) (*statsResponse, error) {
	check := cs.createCheck("collect_stats")
	defer cs.monitorCheck(time.Now(), check)

	body, err := client.get("/stats")
	if err != nil {
		check.fail(fmt.Sprintf("An error occured fetching gossip. %s", err))
		return nil, err
	}

	r, err := toStatsResponse(body)
	if err != nil {
		check.fail(fmt.Sprintf("An error occured parsing gossip. %s", err))
		return nil, err
	}

	return r, nil
}

func (cs *checkSet) doSysCPUCheck(r *statsResponse) {
	check := cs.createCheck("sys_cpu")
	defer cs.monitorCheck(time.Now(), check)
	expected := 90.0 // TODO: Make configurable

	cpu := math.Round(r.Sys.CPU*100) / 100

	check.Data = cpu
	check.Output = fmt.Sprintf("%.2f%% cpu in use by system", cpu)
	if cpu > expected {
		check.warn(fmt.Sprintf("System is using a lot of cpu! %.2f%%.", cpu))
	}
}

func (cs *checkSet) doSysMemoryCheck(r *statsResponse) {
	check := cs.createCheck("sys_mem")
	defer cs.monitorCheck(time.Now(), check)
	expected := 200 // TODO: Make configurable

	freeMB := int((r.Sys.FreeMemory / 1000) / 1000)

	check.Data = freeMB
	check.Output = fmt.Sprintf("%dMB system memory free", freeMB)
	if freeMB < expected {
		check.warn(fmt.Sprintf("Free system memory is low! %dMB free.", freeMB))
	}
}

func (cs *checkSet) doProcCPUCheck(r *statsResponse) {
	check := cs.createCheck("proc_cpu")
	defer cs.monitorCheck(time.Now(), check)
	expected := 90.0 // TODO: Make configurable

	cpu := math.Round(r.Proc.CPU*100) / 100

	check.Data = cpu
	check.Output = fmt.Sprintf("%.2f%% cpu in use by process", cpu)
	if cpu > expected {
		check.warn(fmt.Sprintf("Process is using a lot of cpu! %.2f%%.", cpu))
	}
}

func (cs *checkSet) doProcMemoryCheck(r *statsResponse) {
	check := cs.createCheck("proc_mem")
	defer cs.monitorCheck(time.Now(), check)
	expected := 1000 // TODO: Make configurable

	usedMB := int((r.Proc.Memory / 1000) / 1000)

	check.Data = usedMB
	check.Output = fmt.Sprintf("%dMB memory used by process", usedMB)
	if usedMB > expected {
		check.warn(fmt.Sprintf("Process is using a a lot of memory! %dMB used.", usedMB))
	}
}
