package collector

import (
	"strconv"
	"strings"

	"github.com/mtulio/statuscake-exporter/stk"
	"github.com/prometheus/client_golang/prometheus"
)

type stkTestCollector struct {
	stkTestUp     *prometheus.Desc
	stkTestUptime *prometheus.Desc
	stkTestPerf   *prometheus.Desc
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
			"StatusCake Test Status",
			[]string{"name"}, nil,
		),
		stkTestUptime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkTestCollectorSubsystem, "uptime"),
			"StatusCake Test Uptime from the last 7 day",
			[]string{"name"}, nil,
		),
		stkTestPerf: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkTestCollectorSubsystem, "performance_ms"),
			"StatusCake Test performance data",
			[]string{"name", "location", "status"}, nil,
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
	if len(tests) < 1 {
		return nil
	}
	for t := range tests {
		test := tests[t]
		testStatus := 0
		if strings.ToLower(test.Status) == "up" {
			testStatus = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkTestUp,
			prometheus.GaugeValue,
			float64(testStatus),
			test.WebsiteName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.stkTestUptime,
			prometheus.GaugeValue,
			float64(test.Uptime),
			test.WebsiteName,
		)
		if len(test.PerformanceData) > 0 {
			for p := range test.PerformanceData {
				ch <- prometheus.MustNewConstMetric(
					c.stkTestPerf,
					prometheus.GaugeValue,
					float64(test.PerformanceData[p].Performance),
					test.WebsiteName,
					test.PerformanceData[p].Location,
					strconv.Itoa(test.PerformanceData[p].Status),
				)
			}
		}
	}

	return nil
}
