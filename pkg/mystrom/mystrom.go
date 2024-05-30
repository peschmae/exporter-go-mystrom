package mystrom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const namespace = "mystrom"

// 5 second timeout, might need to be increased
const reqTimeout = time.Second * 5

type switchReport struct {
	Power       float64 `json:"power"`
	WattPerSec  float64 `json:"Ws"`
	Relay       bool    `json:"relay"`
	Temperature float64 `json:"temperature"`
}

type switchInfo struct {
	Version   string  `json:"version"`
	Mac       string  `json:"mac"`
	SwType    float64 `json:"type"`
	SSID      string  `json:"ssid"`
	Static    bool    `json:"static"`
	Connected bool    `json:"connected"`
}

// Exporter --
type Exporter struct {
	myStromSwitchIp string
	switchType      float64
}

// NewExporter --
func NewExporter(switchIP string) *Exporter {
	return &Exporter{
		myStromSwitchIp: switchIP,
	}
}

// Scrape --
func (e *Exporter) Scrape() (prometheus.Gatherer, error) {
	reg := prometheus.NewRegistry()

	// --
	bodyInfo, err := e.fetchData("/api/v1/info")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to target: %v", err.Error())
	}

	info := switchInfo{}
	err = json.Unmarshal(bodyInfo, &info)
	if err != nil {
		return reg, fmt.Errorf("unable to decode switchReport: %v", err.Error())
	}
	log.Debugf("info: %#v", info)
	e.switchType = info.SwType

	if err := registerInfoMetrics(reg, info, e.myStromSwitchIp); err != nil {
		return nil, fmt.Errorf("failed to register metrics : %v", err.Error())
	}

	// --
	bodyData, err := e.fetchData("/report")
	if err != nil {
		return reg, fmt.Errorf("unable to fetch switchReport: %v", err.Error())
	}

	report := switchReport{}
	err = json.Unmarshal(bodyData, &report)
	if err != nil {
		return reg, fmt.Errorf("unable to decode switchReport: %v", err.Error())
	}
	log.Debugf("report: %#v", report)

	if err := registerMetrics(reg, report, e.myStromSwitchIp, e.switchType); err != nil {
		return nil, fmt.Errorf("failed to register metrics : %v", err.Error())
	}

	return reg, nil
}

// fetchData -- get the data from the switch under the given path
func (e *Exporter) fetchData(urlpath string) ([]byte, error) {
	url := "http://" + e.myStromSwitchIp + urlpath

	switchClient := http.Client{
		Timeout: reqTimeout,
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "myStrom-exporter")

	res, err := switchClient.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to connect to target: %v", err.Error())
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to read body: %v", err.Error())
	}

	return body, nil
}

// registerMetrics --
func registerMetrics(reg prometheus.Registerer, data switchReport, target string, st float64) error {

	// --
	collectorRelay := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "relay",
			Help:      "The current state of the relay (whether or not the relay is currently turned on)",
		},
		[]string{"instance"})

	if err := reg.Register(collectorRelay); err != nil {
		return fmt.Errorf("failed to register metric %v: %v", "relay", err.Error())
	}

	if data.Relay {
		collectorRelay.WithLabelValues(target).Set(1)
	} else {
		collectorRelay.WithLabelValues(target).Set(0)
	}

	if st != 114 {
		// --
		collectorPower := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "power",
				Help:      "The current power consumed by devices attached to the switch",
			},
			[]string{"instance"})

		if err := reg.Register(collectorPower); err != nil {
			return fmt.Errorf("failed to register metric %v: %v", "power", err.Error())
		}

		collectorPower.WithLabelValues(target).Set(data.Power)

		// --
		collectorAveragePower := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "average_power",
				Help:      "The average power since the last call. For continuous consumption measurements.",
			},
			[]string{"instance"})

		if err := reg.Register(collectorAveragePower); err != nil {
			return fmt.Errorf("failed to register metric %v: %v", "average_power", err)
		}

		collectorAveragePower.WithLabelValues(target).Set(data.WattPerSec)

		// --
		collectorTemperature := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "temperature",
				Help:      "The currently measured temperature by the switch. (Might initially be wrong, but will automatically correct itself over the span of a few hours)",
			},
			[]string{"instance"})

		if err := reg.Register(collectorTemperature); err != nil {
			return fmt.Errorf("failed to register metric %v: %v", "temperature", err.Error())
		}

		collectorTemperature.WithLabelValues(target).Set(data.Temperature)

	}

	return nil
}

// registerMetrics --
func registerInfoMetrics(reg prometheus.Registerer, data switchInfo, target string) error {

	// --
	collectorInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "info",
			Help:      "general information about the device",
		},
		[]string{"instance", "version", "mac", "type", "ssid"})

	if err := reg.Register(collectorInfo); err != nil {
		return fmt.Errorf("failed to register metric %v: %v", "info", err.Error())
	}

	collectorInfo.WithLabelValues(target, data.Version, data.Mac, fmt.Sprintf("%v", data.SwType), data.SSID).Set(1)

	return nil
}
