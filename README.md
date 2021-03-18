# myStrom prometheus exporter

Export [myStrom WiFi Switch](https://mystrom.ch/de/wifi-switch-ch/) report
statistics to [Prometheus](https://prometheus.io).

Metrics are retrieved using the [switch REST api](https://api.mystrom.ch/).
This has only be testet using a WiFi switch with firmware 3.82.60, but should be 
backwards compatible with all 3.x firmwares.

To run it:
```bash
$ go build
$ ./mystrom-exporter [flags]
```

## Exported Metrics
| Metric | Description |
| ------ | ------- |
| mystrom_up | Was the last REST api call to the switch successful |
| mystrom_report_watt_per_sec | The average of energy consumed per second from last call this request |
| mystrom_report_temperatur  | The currently measured temperature by the switch. (Might initially be wrong, but will automatically correct itself over the span of a few hours) |
| mystrom_report_relay | The current state of the relay (wether or not the relay is currently turned on) |
| mystrom_report_power  | The current power consumed by devices attached to the switch |

## Flags
```bash
$ ./mystrom-exporter --help
```
| Flag | Description | Default |
| ---- | ----------- | ------- |
| switch.ip-address | IP address of the switch you try to monitor | `` |
| web.listen-address | Address to listen on | `:9452` |
| web.metrics-path | Path under which to expose metrics | `/metrics` |


## License
MIT License, See the included LICENSE file for terms and conditions.