"""Module contianing functions used to manage prometheus
counters in the Hermes client"""

import logging
import json
from typing import Dict

from pydantic import ValidationError

from hermes import hermes
from hermes.models import PrometheusCounter
from hermes.hermes import push_udp_packet, set_hermes_config
from hermes.exceptions import InvalidMetricException

LOGGER = logging.getLogger('hermes.counters')


def increment_counter(metric_name: str, labels: Dict[str, str]):
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
    if None in (hermes.HERMES_HOST, hermes.HERMES_PORT):
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

def counter_wrapper(metric_name: str, labels: dict, pre_execution: bool = True):
    """Decorator used to wrap a function in a
    counter function. The counter is incremented
    on function execution

    Arguments:
        metric_name: str name of metric on Hermes Server
        labels: dict labels and their corresponding string
            values
        pre_execution: bool counter is incremented before
            function execution if true, else after
    """
    def wrapper(func: object):
        def make_wrapper(*args, **kwargs):
            # increment counter before job if specified
            if pre_execution:
                increment_counter(metric_name, labels)
                return func(*args, **kwargs)
            # else, execute job first and increment counter
            else:
                results = func(*args, **kwargs)
                increment_counter(metric_name, labels)
                return results
        return make_wrapper
    return wrapper
