# netdisco-exporter

## Prometheus Exporter for metrics from Netdisco

[![Go Report Card](https://goreportcard.com/badge/github.com/cbutera-sqsp/netdisco-exporter)](https://goreportcard.com/report/github.com/cbutera-sqsp/netdisco-exporter)

### Docker
```
docker run --restart unless-stopped -d -p 8080:8080 -e NETDISCO_HOST -e NETDISCO_USERNAME -e NETDISCO_PASSWORD --name netdisco-exporter cbutera90/netdisco-exporter
```

## Metrics
The following metrics are supported by now:
- Netdisco API status
- last discover
- last arpnip
- last macsuck

### Docker image
https://hub.docker.com/r/cbutera90/netdisco-exporter 

## License
(c) Licensed under [MIT](LICENSE) license.

## Prometheus
see https://prometheus.io/

## Netdisco
see https://github.com/netdisco/netdisco 