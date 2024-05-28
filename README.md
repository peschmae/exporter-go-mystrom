# myStrom prometheus exporter

Export [myStrom WiFi Switch](https://mystrom.ch/de/wifi-switch-ch/) report
statistics to [Prometheus](https://prometheus.io).

Metrics are retrieved using the [switch REST api](https://api.mystrom.ch/).
This has only been tested using a Wi-Fi switch with firmware 3.82.60, but should be backwards compatible with all 3.x
firmwares. It was also tested on the 4.0.7 firmware.

To run it:

```bash
$ go build
$ ./mystrom-exporter [flags]
```

## Build instructions

The package uses `stringer` to generate `String()` methods on structs, to build the package you need to
install `stringer` through `gotools`.

```bash
$ go install golang.org/x/tools/cmd/stringer@latest
# optional, should also be triggered by go build
$ go generate ./...
```

## Exported Metrics

| Metric                  | Description                                                                                                                                      |
|-------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| `mystrom_info`          | General information about the device (in the labels)                                                                                             |
| `mystrom_relay`         | The current state of the relay (whether or not the relay is currently turned on)                                                                 |
| `mystrom_power`         | The current power consumed by devices attached to the switch                                                                                     |
| `mystrom_average_power` | The average power since the last call. For continuous consumption measurements.                                                                  |
| `mystrom_temperature`   | The currently measured temperature by the switch. (Might initially be wrong, but will automatically correct itself over the span of a few hours) |

## Flags

```bash
$ ./mystrom-exporter --help
```

| Flag                 | Description                                                                          | Default    |
|----------------------|--------------------------------------------------------------------------------------|------------|
| `web.listen-address` | Address to listen on                                                                 | `:9452`    |
| `web.metrics-path`   | Path under which to expose exporters own metrics                                     | `/metrics` |
| `web.device-path`    | Path under which the metrics of the devices are fetched, requires `target` parameter | `/device`  |

## Prometheus configuration

An enhancement has been made to have only one exporter which can scrape multiple devices.
This is configured in Prometheus as follows assuming we have 4 mystrom devices and the exporter is running locally on
the same machine as the Prometheus.

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
     - source_labels: [ __address__ ]
       target_label: __param_target
     - target_label: __address__
       replacement: 127.0.0.1:9452
```

## Supported architectures

Using the make file, you can easily build for the following architectures, those can also be considered the tested ones:

| OS    | Arch  |
|-------|-------|
| Linux | amd64 |
| Linux | arm64 |
| Linux | arm   |
| Mac   | amd64 |
| Mac   | arm64 |

Since go is cross compatible with windows, and mac arm as well, you should be able to build the binary for those as
well, but they aren't tested.
The docker image is only built & tested for amd64.

## Packages

Packages are built automatically on release, and container images on push to the main branch.

Take a look at the `Releases` or `Packages` tabs on GitHub.

### Container images

There is a multiplatform build available
here https://github.com/peschmae/exporter-go-mystrom/pkgs/container/exporter-go-mystrom

```
docker pull ghcr.io/peschmae/exporter-go-mystrom:latest
```

## License

MIT License, See the included LICENSE file for terms and conditions.