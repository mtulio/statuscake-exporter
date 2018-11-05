package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/mtulio/statuscake-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readVersion() {
	dat, err := ioutil.ReadFile("./VERSION")
	check(err)
	config.version = string(dat)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	prometheus.MustRegister(version.NewCollector("statuscake_exporter"))

	flagListenAddress := flag.String("web.listen-address", config.listenAddress, "Address on which to expose metrics and web interface.")
	flagMetricsPath := flag.String("web.telemetry-path", config.metricsPath, "Path under which to expose metrics.")
	flagStkUsername := flag.String("stk.username", "", "StatusCake API's Username.")
	flagStkApikey := flag.String("stk.apikey", "", "StatusCake API's Apikey.")
	flagStkTags := flag.String("stk.tags", "", "StatusCake Filter Tags separated by comma.")
	flagVersion := flag.Bool("v", false, "prints current version")
	flag.Usage = usage
	flag.Parse()

	readVersion()

	if *flagVersion {
		fmt.Println(config.version)
		os.Exit(0)
	}

	if *flagListenAddress != config.listenAddress {
		config.listenAddress = *flagListenAddress
	}

	if *flagMetricsPath != config.metricsPath {
		config.metricsPath = *flagMetricsPath
	}

	if *flagStkUsername == "" {
		log.Errorln("StatusCake API must be provided.")
		os.Exit(1)
	} else {
		config.StkUsername = *flagStkUsername
	}

	if *flagStkApikey == "" {
		log.Errorln("StatusCake API APIKEY must be provided.")
		os.Exit(1)
	} else {
		config.StkApikey = *flagStkApikey
	}

	if *flagStkTags != "" {
		Stk.Tags = *flagStkTags
	}
}

// Main Prometheus handler
func handler(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	log.Debugln("collect query:", filters)

	mc, err := collector.NewMasterCollector(filters...)
	if err != nil {
		log.Warnln("Couldn't create", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Couldn't create %s", err)))
		return
	}

	registry := prometheus.NewRegistry()
	err = registry.Register(mc)
	if err != nil {
		log.Errorln("Couldn't register collector:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Couldn't register collector: %s", err)))
		return
	}

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(gatherers,
			promhttp.HandlerOpts{
				// ErrorLog:      log.NewErrorLogger(),
				ErrorHandling: promhttp.ContinueOnError,
			}),
	)
	h.ServeHTTP(w, r)
}

func main() {

	log.Infoln("Starting exporter ", config.version)
	log.Infoln("Initializing Status Cake client...", initClient())

	// Instance master collector that will keep all subsystems
	// This instance is only used to check collector creation and logging.
	mc, err := collector.NewMasterCollector()
	if err != nil {
		log.Fatalf("Couldn't create collector: %s", err)
	}
	log.Infof("Enabled collectors:")
	collectors := []string{}
	for n := range mc.Collectors {
		collectors = append(collectors, n)
	}
	sort.Strings(collectors)
	for _, n := range collectors {
		log.Infof(" - %s", n)
	}

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.HandleFunc(config.metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>StatusCake Exporter</title></head>
			<body>
			<h1>StatusCake Exporter</h1>
			<p><a href="` + config.metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Beginning to serve on port " + config.listenAddress)
	log.Fatal(http.ListenAndServe(config.listenAddress, nil))
}
