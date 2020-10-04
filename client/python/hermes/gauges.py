"""Module contianing functions used to manage prometheus
gauges in the Hermes client"""

import logging
import json
from typing import Dict
from contextlib import contextmanager

from pydantic import ValidationError

from hermes.models import PrometheusGauge
from hermes.hermes import HERMES_HOST, HERMES_PORT
from hermes.hermes import push_udp_packet, set_hermes_config
from hermes.exceptions import InvalidMetricException

LOGGER = logging.getLogger('hermes.gauges')


def set_gauge(metric_name: str, value: float, labels: Dict[str, str]):
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
        payload = {'value': value, 'labels': labels, 'operation': 'set'}
        gauge = PrometheusGauge(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(gauge.json()))
    except ValidationError:
        LOGGER.exception('received invalid gauge configuration')
        raise InvalidMetricException

def increment_gauge(metric_name: str, labels: Dict[str, str]):
    """Function used to push new UDP packet for
    Gauge instance defined on Hermes Server. Inputs
    are converted to Pydantic models for validation
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
        LOGGER.debug('incrementing gauge metric \'%s\'', metric_name)
        payload = {'labels': labels, 'operation': 'increment'}
        gauge = PrometheusGauge(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(gauge.json()))
    except ValidationError:
        LOGGER.exception('received invalid gauge configuration')
        raise InvalidMetricException

def decrement_gauge(metric_name: str, labels: Dict[str, str]):
    """Function used to push new UDP packet for
    Gauge instance defined on Hermes Server. Inputs
    are converted to Pydantic models for validation
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
        LOGGER.debug('decrementing gauge metric \'%s\'', metric_name)
        payload = {'labels': labels, 'operation': 'decrement'}
        gauge = PrometheusGauge(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(gauge.json()))
    except ValidationError:
        LOGGER.exception('received invalid gauge configuration')
        raise InvalidMetricException

@contextmanager
def hermes_gauge(metric_name: str, labels: dict):
    """context manager used to incremnet a
    gauge while a particular operation is
    performed. Once the operation is completed,
    the gauge is dercemented

    Arguments:
        metric_name: str name of metric on Hermes Server
        labels: dict labels and their corresponding string
            values
    """
    try:
        yield increment_gauge(metric_name, labels)
    finally:
        decrement_gauge(metric_name, labels)

def gauge_wrapper(metric_name: str, labels: dict):
    """Decorator used to wrap a function in a
    gauge function. The gauge is incremented
    on function execution and decremented on
    function finish

    Arguments:
        metric_name: str name of metric on Hermes Server
        labels: dict labels and their corresponding string
            values
    """
    def wrapper(func: object):
        def make_wrapper(*args, **kwargs):
            # increment gauge and execute function
            increment_gauge(metric_name, labels)
            results = func(*args, **kwargs)
            # decrement gauge and return results
            decrement_gauge(metric_name, labels)
            return results
        return make_wrapper
    return wrapper


