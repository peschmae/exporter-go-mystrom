package mystrom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "mystrom"

type switchReport struct {
	Power       float64 `json:"power"`
	WattPerSec  float64 `json:"Ws"`
	Relay       bool    `json:relay`
	Temperature float64 `json:"temperature`
}

type Exporter struct {
	myStromSwitchIp string
}

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last myStrom query successful.",
		nil, nil,
	)
	myStromPower = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "report_power"),
		"The current power consumed by devices attached to the switch",
		nil, nil,
	)

	myStromRelay = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "report_relay"),
		"The current state of the relay (wether or not the relay is currently turned on)",
		nil, nil,
	)

	myStromTemperature = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "report_temperatur"),
		"The currently measured temperature by the switch. (Might initially be wrong, but will automatically correct itself over the span of a few hours)",
		nil, nil,
	)

	myStromWattPerSec = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "report_watt_per_sec"),
		"The average of energy consumed per second from last call this request",
		nil, nil,
	)
)

func NewExporter() *Exporter {
	return &Exporter{
		// myStromSwitchIp: myStromSwitchIp,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- myStromPower
	ch <- myStromRelay
	ch <- myStromTemperature
	ch <- myStromWattPerSec
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	e.FetchSwitchMetrics(e.myStromSwitchIp, ch)
}

func (e *Exporter) FetchSwitchMetrics(switchIP string, ch chan<- prometheus.Metric) {

	url := "http://" + switchIP + "/report"

	switchClient := http.Client{
		Timeout: time.Second * 5, // 3 second timeout, might need to be increased
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
	}

	req.Header.Set("User-Agent", "myStrom-exporter")

	res, getErr := switchClient.Do(req)
	if getErr != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)

	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
	}

	report := switchReport{}
	err = json.Unmarshal(body, &report)
	if err != nil {
		fmt.Println(err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		myStromPower, prometheus.GaugeValue, report.Power,
	)

	if report.Relay {
		ch <- prometheus.MustNewConstMetric(
			myStromRelay, prometheus.GaugeValue, 1,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			myStromRelay, prometheus.GaugeValue, 0,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		myStromWattPerSec, prometheus.GaugeValue, report.WattPerSec,
	)

	ch <- prometheus.MustNewConstMetric(
		myStromTemperature, prometheus.GaugeValue, report.Temperature,
	)

}

func (e *Exporter) Scrape(targetAddress string) (prometheus.Gatherer, error) {
	reg := prometheus.NewRegistry()

	return reg, nil
}
