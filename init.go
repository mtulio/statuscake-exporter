package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mtulio/statuscake-exporter/collector"
	stk "github.com/mtulio/statuscake-exporter/statusCake"
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

	flagStkUsername := flag.String("stk.username", "", "StatusCake API's Username.")
	flagStkApikey := flag.String("stk.apikey", "", "StatusCake API's Apikey.")
	flagStkTags := flag.String("stk.tags", "", "StatusCake Filter Tags separated by comma.")

	flagVersion := flag.Bool("v", false, "prints current version")
	flag.Usage = usage
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
		config.StkTags = *flagStkTags
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

	log.Info("Initializing Status Cake client...")

	stkAPI, err = stk.NewStkAPI(config.StkUsername, config.StkApikey)
	if err != nil {
		log.Errorln("Init StatusCake API: ", err)
		return err
	}
	log.Infoln("Success")

	err = stkAPI.GatherAll()
	if err != nil {
		log.Warnln("Init Prom: Couldn't create collector: ", err)
		return err
	}
	if config.StkTags != "" {
		stkAPI.SetConfigTags(config.StkTags)
	}

	return nil
}
