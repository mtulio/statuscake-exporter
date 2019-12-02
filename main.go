package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

func init() {
	var err error
	err = nil

	// Flags
	err = initFlags()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	// StatusCake API
	err = initStkAPI()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	// Prometheus
	initPrometheusApp()
	err = initPrometheus()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

// Main Prometheus handler
func handler(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	log.Debugln("collect query:", filters)

	if len(filters) > 0 {
		err := initPrometheus(filters...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Couldn't create %s", err)))
			return
		}
	}

	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.InstrumentMetricHandler(
		prom.Registry,
		promhttp.HandlerFor(prom.Gatherers,
			promhttp.HandlerOpts{
				ErrorHandling: promhttp.ContinueOnError,
			}),
	)
	h.ServeHTTP(w, r)
}

func main() {

	log.Infoln("Starting exporter ", version.Version)

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.HandleFunc(config.metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>"` + exporterDescription + `"</title></head>
			<body>
			<h1>"` + exporterDescription + `"</h1>
			<p><a href="` + config.metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Beginning to serve on port " + config.listenAddress)
	log.Fatal(http.ListenAndServe(config.listenAddress, nil))
}
