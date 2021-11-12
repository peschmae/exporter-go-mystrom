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

## Supported architectures
Using the make file, you can easily build for the following architectures, those can also be considered the tested ones:
| OS | Arch |
| -- | ---- |
| Linux | amd64 |
| Linux | arm64 |
| Linux | arm |
| Mac | amd64 |

Since go is cross compatible with windows, and mac arm as well, you should be able to build the binary for those as well, but they aren't tested.  
The docker image is only built & tested for amd64.

## License
MIT License, See the included LICENSE file for terms and conditions.