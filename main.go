package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"

	"mystrom-exporter/pkg/version"
)

type switchReport struct {
	Power       float64 `json:"power"`
	WattPerSec  float64 `json:"Ws"`
	Relay       bool    `json:"relay"`
	Temperature float64 `json:"temperature"`
}

const namespace = "mystrom"

var (
	listenAddress = flag.String("web.listen-address", ":9452",
		"Address to listen on")
	metricsPath = flag.String("web.metrics-path", "/metrics",
		"Path under which to expose metrics")
	switchIP = flag.String("switch.ip-address", "",
		"IP address of the switch you try to monitor")
	showVersion = flag.Bool("version", false,
		"Show version information.")

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

type Exporter struct {
	myStromSwitchIp string
}

func NewExporter(myStromSwitchIp string) *Exporter {
	return &Exporter{
		myStromSwitchIp: myStromSwitchIp,
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

	report, err := FetchReport(switchIP)

	if err != nil {
		log.Infof("Error occured, while fetching metrics: %s", err)
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

func FetchReport(switchIP string) (*switchReport, error) {
	log.Infof("Trying to connect to switch at: %s", switchIP)
	url := "http://" + switchIP + "/report"

	switchClient := http.Client{
		Timeout: time.Second * 5, // 3 second timeout, might need to be increased
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "myStrom-exporter")

	res, getErr := switchClient.Do(req)
	if getErr != nil {
		log.Infof("Error while trying to connect to switch: %s", getErr)
		return nil, getErr

	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Infof("Error while reading body: %s", readErr)
		return nil, readErr
	}

	report := switchReport{}
	err = json.Unmarshal(body, &report)
	if err != nil {
		log.Infof("Error while unmarshaling report: %s", err)
		return nil, err
	}

	return &report, nil
}

func main() {

	flag.Parse()

	// Show version information
	if *showVersion {
		v, err := version.Print("mystrom_exporter")
		if err != nil {
			log.Fatalf("Failed to print version information: %#v", err)
		}

		fmt.Fprintln(os.Stdout, v)
		os.Exit(0)
	}

	if *switchIP == "" {
		flag.Usage()
		fmt.Println("\nNo switch.ip-address provided")
		os.Exit(1)
	}

	// make the build information is available through a metric
	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "scripts",
			Name:      "build_info",
			Help:      "A metric with a constant '1' value labeled by build information.",
		},
		[]string{"version", "revision", "branch", "goversion", "builddate", "builduser"},
	)
	buildInfo.WithLabelValues(version.Version, version.Revision, version.Branch, version.GoVersion, version.BuildDate, version.BuildUser).Set(1)

	exporter := NewExporter(*switchIP)
	prometheus.MustRegister(exporter, buildInfo)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>myStrom switch report Exporter</title></head>
             <body>
             <h1>myStrom Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	_, err := FetchReport(*switchIP)
	if err != nil {
		log.Fatalf("Switch at address %s couldn't be reached. Ensure it is reachable before starting the exporter", *switchIP)
	}

	log.Infoln("Listening on address " + *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
