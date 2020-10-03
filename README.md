# Hermes: Push Gateway for Prometheus Metrics

UDP Push gateway for prometheus metrics

The `Hermes` server acts as a push-gateway to allow worker processes to
push Prometheus Metrics to an intermediate metric store, which is then
scraped by some Prometheus instance. This avoids the need to run a
Prometheus server instance within each worker process.

Communicate is handled entirely via UDP to provide a low-overhead interface
that will not stop/break any worker processes if the `Hermes` server stops
working (as a HTTP/TCP connection would). The metrics that can be pushed to a
`Hermes` server instance are defined in a JSON file that is mounted into the
container of the service. The JSON file contains definitions for Gauges and
Counters, which are then created at the server start time. JSON payloads
containing the metric name and JSON payload can then be pushed to the `Hermes`
server within the client code via the UDP interface (see `example/send_metric.py`)
for an example of how to push UDP JSON packets in Python

A typical JSON configuration can be defined as follows

```json
{
    "service_name": "testing-service",
    "gauges": [
        {
            "metric_name": "sample_gauge",
            "metric_description": "gauges are awesome",
            "labels": ["label_1", "label_2"]
        }
    ],
    "counters": [
        {
            "metric_name": "sample_counter",
            "metric_description": "counters are awesome too",
            "labels": ["label_1"]
        }
    ]
}
```

Note that currently only one service name is supported per Hermes instance, as Hermes is designed
to be a side-cart application for applications deployed on Docker Swarm, Kubernetes and similar
container orchestration platforms.

The UDP interface by default listens on port `7789`, which can be configured in the environment
variables of the container, while the `Prometheus` interface listens on port `8080`. The UDP packets
send to the Hermes server must have the following format

### Counters

```json
{
    "metric_name": "sample_counter",
    "payload": {
        "labels": {
            "label_1": "testing label 1"
        }
    }
}
```

### Gauges

```json
{
    "metric_name": "gauge",
    "payload": {
        "labels": {
            "label_1": "testing label 1",
            "label_2": "testing label 2"
        },
        "value": 65.4
    }
}
```

Note that the labels defined in the JSON packets must match the labels defined in the
`Hermes` configuration file