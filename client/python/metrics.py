"""Module contianing python client functions for Hermes"""

import logging
import socket
import json
from typing import Dict

from pydantic import BaseModel, ValidationError

SOCKET = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
LOGGER = logging.getLogger(__name__)


class InvalidMetricException(Exception):
    """Exception raised when invalid metrics
    are sent to the Hermes client"""

class GaugePayload(BaseModel):
    """Dataclass containing values for
    Gauge payload"""
    value: float
    labels: Dict[str, str]

class PrometheusGauge(BaseModel):
    """Dataclass containing values for
    a Prometheus Gauge instance"""
    metric_name: str
    payload: GaugePayload

class CounterPayload(BaseModel):
    """Dataclass containing values for
    a Prometheus Counter instance"""
    labels: Dict[str, str]

class PrometheusCounter(BaseModel):
    """Dataclass containing values for
    a Prometheus counter instance"""
    metric_name: str
    payload: CounterPayload

def push_udp_packet(payload: dict, host: str = 'localhost', port: int = 7789):
    """Function used to send UDP packet to hermes server"""
    payload = json.dumps(payload)
    SOCKET.sendto(bytes(payload, encoding='utf-8'), (host, port))

def push_gauge(metric_name: str, value: float, labels: Dict[str, str]):
    """Function used to push new UDP packet for
    Gauge instance defined on Hermes Server

    Arguments:
        metric_name: str name of metric on Hermes Server
        value: float value to set gauge to
    """
    try:
        payload = {'value': value, 'labels': labels}
        gauge = PrometheusGauge(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(gauge.json()))
    except ValidationError:
        LOGGER.exception('received invalid gauge configuration')
        raise InvalidMetricException

def push_counter(metric_name: str, labels: Dict[str, str]):
    """Function used to push new UDP packet for
    counter instance defined on Hermes Server

    Arguments:
        metric_name: str name of metric on Hermes Server
    """
    try:
        payload = {'labels': labels}
        counter = PrometheusCounter(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(counter.json()))
    except ValidationError:
        LOGGER.exception('received invalid gauge configuration')
        raise InvalidMetricException

if __name__ == '__main__':

    push_counter('sample_counter', {'label_1': 'testing python client'})
    push_gauge('sample_gauge', 65.6, {'label_1': 'testing python client', 'label_2': 'testing python client again'})