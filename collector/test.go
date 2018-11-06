package collector

import (
	"strings"

	stk "github.com/mtulio/statuscake-exporter/statusCake"
	"github.com/prometheus/client_golang/prometheus"
)

type stkTestCollector struct {
	stkTestUp     *prometheus.Desc
	stkTestUptime *prometheus.Desc
	StkAPI        *stk.StkAPI
}

const (
	stkTestCollectorSubsystem = "test"
)

func init() {
	registerCollector("test", defaultEnabled, NewStkTestCollector)
}

//NewStkTestCollector is a Status Cake Test Collector
func NewStkTestCollector() (Collector, error) {
	return &stkTestCollector{
		stkTestUp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkTestCollectorSubsystem, "up"),
			"Status Cake test Status",
			[]string{"name"}, nil,
		),
		stkTestUptime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkTestCollectorSubsystem, "uptime"),
			"Status Cake test Uptime from the last 7 day",
			[]string{"name"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes related metrics
func (c *stkTestCollector) Update(ch chan<- prometheus.Metric) error {
	if err := c.updateStkTest(ch); err != nil {
		return err
	}
	return nil
}

func (c *stkTestCollector) UpdateConfig(stkAPI *stk.StkAPI) error {
	c.StkAPI = stkAPI
	return nil
}

func (c *stkTestCollector) updateStkTest(ch chan<- prometheus.Metric) error {

	if c.StkAPI == nil {
		return nil
	}
	tests := c.StkAPI.GetTests()
	for t := range tests {
		testStatus := 0
		if strings.ToLower(tests[t].Status) == "up" {
			testStatus = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkTestUp,
			prometheus.GaugeValue,
			float64(testStatus),
			tests[t].WebsiteName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.stkTestUptime,
			prometheus.GaugeValue,
			float64(tests[t].Uptime),
			tests[t].WebsiteName,
		)
	}

	return nil
}
