package collector

import (
	"strconv"
	"time"
	"math"

	"github.com/mtulio/statuscake-exporter/stk"
	"github.com/prometheus/client_golang/prometheus"
)

type stkSSLCollector struct {
	StkAPI              *stk.StkAPI
	stkSslUp            *prometheus.Desc
	stkSslCertScore     *prometheus.Desc
	stkSslCipherScore   *prometheus.Desc
	stkSslCertStatus    *prometheus.Desc
	stkSslValidUntil    *prometheus.Desc
	stkSslAlertReminder *prometheus.Desc
	stkSslLastReminder  *prometheus.Desc
	stkSslAlertExpiry   *prometheus.Desc
	stkSslAlertBroken   *prometheus.Desc
	stkSslAlertMixed    *prometheus.Desc
}

const (
	stkSSLCollectorSubsystem = "ssl"
)

func init() {
	registerCollector("ssl", defaultEnabled, NewStkSSLCollector)
}

//NewStkSSLCollector is a StatusCake SSL Collector
func NewStkSSLCollector() (Collector, error) {
	return &stkSSLCollector{
		stkSslUp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "up"),
			"StatusCake Test SSL Status",
			[]string{"domain"}, nil,
		),
		stkSslCertScore: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "cert_score"),
			"StatusCake SSL Certificate Score",
			[]string{"domain"}, nil,
		),
		stkSslCipherScore: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "cipher_score"),
			"StatusCake SSL Cipher Score",
			[]string{"domain"}, nil,
		),
		stkSslCertStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "cert_status"),
			"StatusCake SSL Cert Status",
			[]string{"domain"}, nil,
		),
		stkSslValidUntil: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "valid_util_sec"),
			"StatusCake SSL Valid Until",
			[]string{"domain"}, nil,
		),
		stkSslAlertReminder: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "alert_reminder"),
			"StatusCake SSL Alert Reminder boolean",
			[]string{"domain"}, nil,
		),
		stkSslLastReminder: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "last_reminder"),
			"StatusCake SSL Last Reminder",
			[]string{"domain"}, nil,
		),
		stkSslAlertExpiry: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "alert_expiry"),
			"StatusCake SSL Alert Expiry boolean",
			[]string{"domain"}, nil,
		),
		stkSslAlertBroken: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "alert_broken"),
			"StatusCake SSL Alert Broken boolean",
			[]string{"domain"}, nil,
		),
		stkSslAlertMixed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, stkSSLCollectorSubsystem, "alert_mixed"),
			"StatusCake SSL Alert Broken boolean",
			[]string{"domain"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes related metrics
func (c *stkSSLCollector) Update(ch chan<- prometheus.Metric) error {
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
		if test.Paused == false {
			testStatus = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkSslUp,
			prometheus.GaugeValue,
			float64(testStatus),
			test.Domain,
		)

		ctscore, _ := strconv.ParseFloat(test.CertScore, 32)
		ch <- prometheus.MustNewConstMetric(
			c.stkSslCertScore,
			prometheus.GaugeValue,
			ctscore,
			test.Domain,
		)

		ciscore, _ := strconv.ParseFloat(test.CipherScore, 32)
		ch <- prometheus.MustNewConstMetric(
			c.stkSslCipherScore,
			prometheus.GaugeValue,
			ciscore,
			test.Domain,
		)

		certStatus := 0
		if test.CertStatus == "CERT_OK" {
			certStatus = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkSslCertStatus,
			prometheus.GaugeValue,
			float64(certStatus),
			test.Domain,
		)

		t, err := time.Parse("2006-01-02 15:04:05", test.ValidUntilUtc)
		if err == nil {
			secUntilExpiry := math.Round(t.Sub(time.Now().UTC()).Seconds())
			ch <- prometheus.MustNewConstMetric(
				c.stkSslValidUntil,
				prometheus.GaugeValue,
				secUntilExpiry,
				test.Domain,
			)
		}

		alert := 0
		if test.AlertReminder {
			alert = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkSslAlertReminder,
			prometheus.GaugeValue,
			float64(alert),
			test.Domain,
		)

		alert = test.LastReminder
		ch <- prometheus.MustNewConstMetric(
			c.stkSslLastReminder,
			prometheus.GaugeValue,
			float64(alert),
			test.Domain,
		)

		alert = 0
		if test.AlertExpiry {
			alert = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkSslAlertExpiry,
			prometheus.GaugeValue,
			float64(alert),
			test.Domain,
		)

		alert = 0
		if test.AlertBroken {
			alert = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkSslAlertBroken,
			prometheus.GaugeValue,
			float64(alert),
			test.Domain,
		)

		alert = 0
		if test.AlertMixed {
			alert = 1
		}
		ch <- prometheus.MustNewConstMetric(
			c.stkSslAlertMixed,
			prometheus.GaugeValue,
			float64(alert),
			test.Domain,
		)
	}

	return nil
}
