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

The following `Dockerfile` illustrates how to use the `Hermes` image

```dockerfile
FROM psauerborn/hermes

COPY ./hermes_config.json ./

ENV CONFIG_FILE_PATH=./hermes_config.json
```

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
    "metric_name": "sample_gauge",
    "payload": {
        "labels": {
            "label_1": "testing label 1",
            "label_2": "testing label 2"
        },
        "operation": "increment"
    }
}
```

```json
{
    "metric_name": "sample_gauge",
    "payload": {
        "labels": {
            "label_1": "testing label 1",
            "label_2": "testing label 2"
        },
        "operation": "decrement"
    }
}
```

```json
{
    "metric_name": "sample_gauge",
    "payload": {
        "labels": {
            "label_1": "testing label 1",
            "label_2": "testing label 2"
        },
        "value": 65.4,
        "operation": "set"
    }
}
```

### Histograms

```json
{
    "metric_name": "sample_histogram",
    "payload": {
        "labels": {
            "label_1": "testing label 1"
        },
        "observation": 65.4
    }
}
```

### Summaries

```json
{
    "metric_name": "sample_summary",
    "payload": {
        "labels": {
            "label_1": "testing label 1"
        },
        "observation": 65.4
    }
}
```

Note that the labels defined in the JSON packets must match the labels defined in the
`Hermes` configuration file

## Python Client Library

`Hermes` has a client library written in python (Go version coming soon). The package can be
installed with

```console
pip install python-hermes
```

Two high-level utility functions are provided, along with context managers and
decorators that can be used to automatically manage gauge and counter interfaces.
See the examples in `example/python` directory for implementation examples

### Counters

```python
import time
from hermes import set_hermes_config, increment_counter, counter_wrapper


def example_job():
    """Example job used to demonstrate counters
    set manually using the provided functions"""
    increment_counter('sample_counter', {'label_1': 'demo label 1'})
    time.sleep(5)

@counter_wrapper('sample_counter', {'label_1': 'demo label 1'})
def example_job_wrapper_pre():
    """Example job used to demonstrate counters
    set with the decorator"""
    print('this counter has been incremented on start')
    time.sleep(5)

@counter_wrapper('sample_counter', {'label_1': 'demo label 1'}, pre_execution=False)
def example_job_wrapper_post():
    """Example job used to demonstrate counters
    set with the decorator"""
    print('this counter will be incremented when im finished')
    time.sleep(5)


if __name__ == '__main__':

    set_hermes_config('localhost', 7789)
    example_job_wrapper_post()
```

### Gauges

```python
from hermes import set_gauge, increment_gauge, decrement_gauge, \
    hermes_gauge, gauge_wrapper, set_hermes_config


def example_job():
    """Example job used to demonstrate gauges
    set manually using the provided functions"""
    increment_gauge('sample_gauge', {'label_1': 'demo label 1', 'label_2': 'demo label 2'})
    time.sleep(5)
    decrement_gauge('sample_gauge', {'label_1': 'demo label 1', 'label_2': 'demo label 2'})

def example_job_context():
    """Example job used to demonstrate gauges set
    with the context manager provided"""
    labels = {'label_1': 'demo label 1', 'label_2': 'demo label 2'}
    with hermes_gauge('sample_gauge', labels=labels):
        print('this gauge has been incremented on start and will decrement when Im done')
        time.sleep(5)

@gauge_wrapper('sample_gauge', {'label_1': 'demo label 1', 'label_2': 'demo label 2'})
def example_job_wrapper():
    """Example job used to demonstrate gauges
    set with the decorator"""
    print('this gauge has been incremented on start and will decrement when Im done')
    time.sleep(5)


if __name__ == '__main__':

    set_hermes_config('localhost', 7789)
    example_job_wrapper()
```

### Histograms

```python
from hermes import observe_histogram

if __name__ == '__main__':

    set_hermes_config('localhost', 7789)

    labels = {'label_1': 'test-label'}
    observe_histogram('sample_histogram', labels=labels, observation=6.54)
```

### Summaries

```python
from hermes import observe_summary

if __name__ == '__main__':

    set_hermes_config('localhost', 7789)

    labels = {'label_1': 'test-label'}
    observe_summary('sample_summary', labels=labels, observation=20.84)
```

## Go Client Library

The `Hermes` package also ships with a `Go` client library. The following code snippet demonstrates
how the Go Client can be used

```go
package main

import (
    "github.com/PSauerborn/hermes/pkg/client"
)

func main() {
    // set configuration for Hermes client and increment sample counter
    client := hermes_client.New("192.168.99.100", 7789)

    // increment counter metric over hermes server
    client.IncrementCounter("sample_counter",
        map[string]string{"label_1": "test-label"})

    // increment gauge
    client.IncrementGauge("sample_gauge",
        map[string]string{"label_1": "test-label", "label_2": "test-label-2"})

    // decrement gauge
    client.DecrementGauge("sample_gauge",
        map[string]string{"label_1": "test-label", "label_2": "test-label-2"})

    // set value of gauge
    client.SetGauge("sample_gauge",
        map[string]string{"label_1": "test-label", "label_2": "test-label-random"}, 54)

    // make observation on histogram
    client.ObserveHistogram("sample_histogram",
        map[string]string{"label_1": "test-label"}, 1233)

    // make observation on summary
    client.ObserveSummary("sample_summary",
        map[string]string{"label_1": "test-label"}, 5)
}
}
```