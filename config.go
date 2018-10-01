package main

import (
	"os"
	"strconv"
	"strings"
)

type emonConfig struct {
	EmonHTTPBindAddress string
	ClusterSize         int
}

var config *emonConfig

func configureEmon() {
	emonHTTPBindAddress := envOrDefault("EMON_HTTP_BIND_ADDRESS", ":8113")
	clusterSize, _ := strconv.Atoi(envOrDefault("EMON_CLUSTER_SIZE", "3"))

	config = &emonConfig{
		EmonHTTPBindAddress: emonHTTPBindAddress,
		ClusterSize:         clusterSize,
	}
}

func envOrDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}
