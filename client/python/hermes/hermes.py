"""Module contianing python client functions for Hermes"""

import logging
import socket
import json
from typing import Dict

from pydantic import BaseModel, ValidationError

from hermes.models import PrometheusGauge, PrometheusCounter
from hermes.exceptions import HermesConfigurationException, InvalidMetricException


SOCKET = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
LOGGER = logging.getLogger('hermes')

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
