package collector

import (
	"fmt"

	"github.com/mtulio/statuscake-exporter/stk"
	"github.com/prometheus/client_golang/prometheus"
)

type stkSSLCollector struct {
	stkTestUp     *prometheus.Desc
	stkTestUptime *prometheus.Desc
	stkTestPerf   *prometheus.Desc
	StkAPI        *stk.StkAPI
}

const (
	stkSSLCollectorSubsystem = "ssl"
)

func init() {
	registerCollector("ssl", defaultEnabled, NewStkSSLCollector)
}

//NewStkSSLCollector is a StatusCake SSL Collector
func NewStkSSLCollector() (Collector, error) {
	return &stkTestCollector{
		stkTestUp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "up"),
			"StatusCake Test Status",
			[]string{"name"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes related metrics
func (c *stkSSLCollector) Update(ch chan<- prometheus.Metric) error {
	fmt.Println("ssl.Update()")
	if err := c.updateStkSSL(ch); err != nil {
		return err
	}
	return nil
}

func (c *stkSSLCollector) UpdateConfig(stkAPI *stk.StkAPI) error {
	c.StkAPI = stkAPI
	return nil
}

func (c *stkSSLCollector) updateStkSSL(ch chan<- prometheus.Metric) error {
	if c.StkAPI == nil {
		return nil
	}
	tests := c.StkAPI.GetTestsSSL()
	if len(tests) < 1 {
		return nil
	}
	for t := range tests {
		test := tests[t]
		testStatus := 0
		fmt.Println(test)
		if test.Paused == false {
			testStatus = 1
		}
		fmt.Println(testStatus)
		fmt.Println(test)
	}

	return nil
}
