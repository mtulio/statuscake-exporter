package main

type globalConf struct {
	listenAddress string
	metricsPath   string
	version       string
	versionCm     string
	versionTag    string
	versionEnv    string
	StkUsername   string
	StkApikey     string
}

var (
	VersionCommit string
	VersionTag    string
	VersionFull   string
	VersionEnv    string
	config        = globalConf{
		":9190",
		"/metrics",
		VersionFull,
		VersionCommit,
		VersionTag,
		VersionEnv,
		"",
		"",
	}
)
