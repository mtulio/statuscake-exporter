package main

type globalConf struct {
	listenAddress string
	metricsPath   string
	version       string
	StkUsername   string
	StkApikey     string
}

var (
	config = globalConf{
		":9190",
		"/metrics",
		"",
		"",
		"",
	}
)
