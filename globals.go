package main

import (
	"github.com/mtulio/statuscake-exporter/collector"
	"github.com/mtulio/statuscake-exporter/stk"
	"github.com/prometheus/client_golang/prometheus"
)

type globalConf struct {
	listenAddress string
	metricsPath   string
	version       string
	versionCm     string
	versionTag    string
	versionEnv    string
	StkUsername   string
	StkApikey     string
	StkTags       string
	StkInterval   int
}

type globalProm struct {
	MC        *collector.MasterCollector
	Registry  *prometheus.Registry
	Gatherers *prometheus.Gatherers
}

const (
	exporterName        = "statuscake_exporter"
	exporterDescription = "StatusCake Exporter"
	default_stkInterval = 300
)

var (
	// VersionCommit is a compiler exporterd var
	VersionCommit string
	VersionTag    string
	VersionFull   string
	VersionEnv    string

	// Global vars
	config = globalConf{
		":9190",
		"/metrics",
		VersionFull,
		VersionCommit,
		VersionTag,
		VersionEnv,
		"",
		"",
		"",
		default_stkInterval,
	}
	stkAPI *stk.StkAPI
	prom   globalProm
)
