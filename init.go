package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mtulio/statuscake-exporter/collector"
	"github.com/mtulio/statuscake-exporter/stk"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

// Flags setup
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func initFlags() error {
	flagListenAddress := flag.String("web.listen-address", config.listenAddress, "Address on which to expose metrics and web interface.")
	flagMetricsPath := flag.String("web.telemetry-path", config.metricsPath, "Path under which to expose metrics.")

	flagStkUsername := flag.String("stk.username", os.Getenv("STATUSCAKE_USER"), "StatusCake API's Username. Default: nil (required)")
	flagStkApikey := flag.String("stk.apikey", os.Getenv("STATUSCAKE_APIKEY"), "StatusCake API's Apikey. Default: nil (required)")
	flagStkTags := flag.String("stk.tags", "", "StatusCake Filter Tags separated by comma. Default: <empty>")
	flagStkInterval := flag.Int("stk.interval", defaultInterval, "StatusCake interval time, in seconds, to gather metrics on API (avoid throtling). Default: 300.")
	flagEnableTests := flag.Bool("stk.enable-tests", true, "Enable Tests module")
	flagEnableSSL := flag.Bool("stk.enable-ssl", true, "Enable SSL module")
	flagSSLFlags := flag.String("stk.ssl-flags", "", "List of flags to expose as metrics sepparated by comma")

	flag.Usage = usage
	flagVersion := flag.Bool("v", false, "prints current version")
	flag.Parse()

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
		log.Errorln("StatusCake API user, or env var STATUSCAKE_USER, must be provided.")
		os.Exit(1)
	} else {
		config.StkUsername = *flagStkUsername
	}

	if *flagStkApikey == "" {
		log.Errorln("StatusCake API APIKEY, or env var STATUSCAKE_APIKEY, must be provided.")
		os.Exit(1)
	} else {
		config.StkApikey = *flagStkApikey
	}

	if *flagStkTags != "" {
		config.StkTags = *flagStkTags
	}

	if *flagStkInterval != defaultInterval {
		config.StkInterval = *flagStkInterval
	}

	config.StkEnableTests = *flagEnableTests
	config.StkEnableSSL = *flagEnableSSL

	if *flagSSLFlags != "" {
		config.StkSSLFlags = *flagSSLFlags
	}

	return nil
}

// Prometheus
func initPrometheusApp() {
	prometheus.MustRegister(version.NewCollector(exporterName))
}

func initPrometheus(filters ...string) error {
	var err error
	err = nil

	prom.MC, err = collector.NewMasterCollector(stkAPI, filters...)
	if err != nil {
		log.Warnln("Init Prom: Couldn't create collector: ", err)
		return err
	}

	prom.Registry = prometheus.NewRegistry()
	err = prom.Registry.Register(prom.MC)
	if err != nil {
		log.Errorln("Init Prom: Couldn't register collector:", err)
		return err
	}

	prom.Gatherers = &prometheus.Gatherers{
		prometheus.DefaultGatherer,
		prom.Registry,
	}
	return nil
}

// StatusCake API setup
func initStkAPI() error {
	var err error
	err = nil

	stkAPI, err = stk.NewStkAPI(config.StkUsername, config.StkApikey)
	if err != nil {
		log.Errorln("Initializing StatusCake API: ", err)
		return err
	}
	log.Info("Initializing StatusCake API client...Success")

	stkAPI.SetWaitInterval(uint32(config.StkInterval))

	err = stkAPI.GatherAll()
	if err != nil {
		log.Warnln("Init Prom: Couldn't create collector: ", err)
		return err
	}

	log.Info("StatusCake collector config:")
	log.Info("- username: ", config.StkUsername)
	log.Info("- interval: ", stkAPI.GetWaitInterval())
	if config.StkTags != "" {
		stkAPI.SetConfigTags(config.StkTags)
		log.Info("- tags: ", stkAPI.GetTags())
	}

	if config.StkSSLFlags != "" {
		for _, s := range strings.Split(config.StkSSLFlags, ",") {
			stkAPI.SetSSLFlag(s)
		}
		log.Info("- SSL flags: ", stkAPI.GetSSLFlags())
	}

	return nil
}
