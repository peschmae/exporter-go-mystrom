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
| web.listen-address | Address to listen on | `:9452` |
| web.metrics-path | Path under which to expose exporters own metrics | `/metrics` |
| web.device-path | Path under which the metrics of the devices are fetched | `/device` |

## Prometheus configuration
A enhancement has been made to have only one exporter which can scrape multiple devices. This is configured in 
Prometheus as follows assuming we have 4 mystrom devices and the exporter is running locally on the smae machine as 
the Prometheus.
```yaml
 - job_name: mystrom
   scrape_interval: 30s
   metrics_path: /device
   honor_labels: true
   static_configs:
   - targets:
     - '192.168.105.11'
     - '192.168.105.12'
     - '192.168.105.13'
     - '192.168.105.14'
   relabel_configs:
     - source_labels: [__address__]
       target_label: __param_target
     - target_label: __address__
       replacement: 127.0.0.1:9452
```

## Supported architectures
Using the make file, you can easily build for the following architectures, those can also be considered the tested ones:
| OS | Arch |
| -- | ---- |
| Linux | amd64 |
| Linux | arm64 |
| Linux | arm |
| Mac | amd64 |
| Mac | arm64 |

Since go is cross compatible with windows, and mac arm as well, you should be able to build the binary for those as well, but they aren't tested.  
The docker image is only built & tested for amd64.

## Packages
Packages are built automatically on release, and container images on push to the main branch.

Take a look at the `Releases` or `Packages` tabs on Github.  

### Container images
There is a multiplatform build available here https://github.com/peschmae/exporter-go-mystrom/pkgs/container/exporter-go-mystrom
```
docker pull ghcr.io/peschmae/exporter-go-mystrom:latest
```

## License
MIT License, See the included LICENSE file for terms and conditions.