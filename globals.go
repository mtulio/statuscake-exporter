package main

import (
	"github.com/mtulio/statuscake-exporter/collector"
	"github.com/mtulio/statuscake-exporter/stk"
	"github.com/prometheus/client_golang/prometheus"
)

type globalConf struct {
	listenAddress      string
	metricsPath        string
	StkUsername        string
	StkApikey          string
	StkTags            string
	StkInterval        int
	Resolution         int
	StkEnableTests     bool
	StkEnableTestsPerf bool
	StkEnableSSL       bool
	StkSSLFlags        string
}

type globalProm struct {
	MC        *collector.MasterCollector
	Registry  *prometheus.Registry
	Gatherers *prometheus.Gatherers
}

const (
	exporterName        = "statuscake_exporter"
	exporterDescription = "StatusCake Exporter"
	defaultInterval     = 300
	defaultResolution   = 30
)

var (
	// Global vars
	config = globalConf{
		":9190",
		"/metrics",
		"",
		"",
		"",
		defaultInterval,
		defaultResolution,
		false,
		false,
		false,
		"",
	}
	stkAPI *stk.StkAPI
	prom   globalProm
)
