package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

func main() {
	log.SetHandler(text.New(os.Stderr))

	configureEmon()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":    r.Method,
			"path":      r.RequestURI,
			"direction": "incoming",
		}).Infof("HTTP %s %s", r.Method, r.RequestURI)

		if strings.HasPrefix(r.RequestURI, "/favicon.ico") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		status := "healthy"
		statusCode := http.StatusOK

		checkSets, code := runHealthchecks()
		switch code {
		case 1:
			status = "warn"
		case 2:
			status = "alert"
			statusCode = http.StatusFailedDependency
		}

		resp := apiResponse{
			Status: status,
			Checks: make(map[string]apiChecks),
		}

		for _, cs := range checkSets {
			if resp.Checks[cs.source] == nil {
				resp.Checks[cs.source] = apiChecks{}
			}
			for _, c := range cs.checks {
				resp.Checks[cs.source][c.Name] = c
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		js, err := json.Marshal(resp)
		if err == nil {
			w.Write(js)
		}

		log.WithFields(log.Fields{
			"status":    statusCode,
			"direction": "outgoing",
		}).Infof("HTTP %s %s - %s", r.Method, r.RequestURI, http.StatusText(statusCode))
	})

	log.WithFields(log.Fields{
		"EmonHTTPBindAddress":    config.EmonHTTPBindAddress,
		"EmonSlowCheckThreshold": config.EmonSlowCheckThreshold,
		"ClusterHTTPEndpoint":    config.ClusterHTTPEndpoint,
		"ClusterSize":            config.ClusterSize,
	}).Infof("emon started")

	srv := &http.Server{
		Addr:         config.EmonHTTPBindAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	srv.ListenAndServe()
}
