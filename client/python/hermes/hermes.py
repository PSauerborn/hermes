"""Module contianing python client functions for Hermes"""

import logging
import socket
import json
from typing import Dict

from pydantic import BaseModel, ValidationError

from hermes.models import PrometheusGauge, PrometheusCounter
from hermes.exceptions import HermesConfigurationException, InvalidMetricException


SOCKET = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
LOGGER = logging.getLogger(__name__)

HERMES_HOST, HERMES_PORT = None, None


def set_hermes_config(host: str = 'localhost', port: int = 7789):
    """Function used to configure the hermes server
    with the correct username and password

    Arguments:
        host: str host of hermes server
        port: int port of hermes server
    """
    global HERMES_HOST, HERMES_PORT
    if not isinstance(port, int) or not isinstance(host, str):
        raise HermesConfigurationException('host and port must be string and int respectively. got \'{}\' \'{}\''.format(type(host), type(port)))
    HERMES_HOST, HERMES_PORT = host, port

def push_udp_packet(payload: dict):
    """Function used to send UDP packet to hermes server"""
    payload = json.dumps(payload)
    SOCKET.sendto(bytes(payload, encoding='utf-8'), (HERMES_HOST, HERMES_PORT))

def push_gauge(metric_name: str, value: float, labels: Dict[str, str]):
    """Function used to push new UDP packet for
    Gauge instance defined on Hermes Server. Inputs
    are converted to Pydantic models for validation
    before they are converted to JSON and pushed
    to the Hermes Server

    Arguments:
        metric_name: str name of metric on Hermes Server
        value: float value to set gauge to
        labels: dict labels and their corresponding string
            values
    """
    # set hermes configuration if not been set before
    if None in (HERMES_HOST, HERMES_PORT):
        LOGGER.warn('hermes configuration not set. setting host to default localhost:7789')
        set_hermes_config()
    try:
        LOGGER.debug('pusing new gauge metric \'%s\' with value %d', metric_name, value)
        payload = {'value': value, 'labels': labels}
        gauge = PrometheusGauge(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(gauge.json()))
    except ValidationError:
        LOGGER.exception('received invalid gauge configuration')
        raise InvalidMetricException

def push_counter(metric_name: str, labels: Dict[str, str]):
    """Function used to push new UDP packet for
    counter instance defined on Hermes Server.
    Inputs are converted to Pydantic models for validation
    before they are converted to JSON and pushed
    to the Hermes Server

    Arguments:
        metric_name: str name of metric on Hermes Server
        labels: dict labels and their corresponding string
            values
    """
    # set hermes configuration if not been set before
    if None in (HERMES_HOST, HERMES_PORT):
        LOGGER.warn('hermes configuration not set. setting host to default localhost:7789')
        set_hermes_config()
    try:
        LOGGER.debug('incrementing counter metric \'%s\'', metric_name)
        payload = {'labels': labels}
        counter = PrometheusCounter(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(counter.json()))
    except ValidationError:
        LOGGER.exception('received invalid gauge configuration')
        raise InvalidMetricException
