package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

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

		status := http.StatusOK

		checkSets, code := checkHealth()
		if code != 0 {
			status = http.StatusFailedDependency
		}

		resp := apiResponse{
			Ok:     true,
			Checks: make(map[string]*check),
		}
		for _, cs := range checkSets {
			for _, c := range cs.checks {
				resp.Checks[c.Name] = c
				if c.Status == statusFailed {
					resp.Ok = false
				}
			}
		}

		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

		log.WithFields(log.Fields{
			"status":    status,
			"direction": "outgoing",
		}).Infof("HTTP %s %s - %s", r.Method, r.RequestURI, http.StatusText(status))
	})

	log.WithFields(log.Fields{
		"EmonHTTPBindAddress": config.EmonHTTPBindAddress,
		"ClusterHTTPEndpoint": config.ClusterHTTPEndpoint,
		"ClusterSize":         config.ClusterSize,
	}).Infof("emon started")
	http.ListenAndServe(config.EmonHTTPBindAddress, nil)
}

type apiResponse struct {
	Ok     bool              `json:"ok"`
	Checks map[string]*check `json:"checks"`
}
