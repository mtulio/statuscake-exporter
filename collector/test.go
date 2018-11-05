package collector

import (
	"log"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type stkTestCollector struct {
	stkTestUp     *prometheus.Desc
	stkTestUptime *prometheus.Desc
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

func (c *stkTestCollector) updateStkTest(ch chan<- prometheus.Metric) error {

	v := url.Values{}
	if Stk.Tags != "" {
		v.Set("tags", Stk.Tags)
	}
	testsWithFilter, err := Stk.Client.Tests().AllWithFilter(v)
	if err != nil {
		log.Fatal(err)
	}
	for t := range testsWithFilter {
		testStatus := 0
		if strings.ToLower(testsWithFilter[t].Status) == "up" {
			testStatus = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkTestUp,
			prometheus.GaugeValue,
			float64(testStatus),
			testsWithFilter[t].WebsiteName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.stkTestUptime,
			prometheus.GaugeValue,
			float64(testsWithFilter[t].Uptime),
			testsWithFilter[t].WebsiteName,
		)
	}

	return nil
}
