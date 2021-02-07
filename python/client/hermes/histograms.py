"""module containing codebase for histogram metrics"""

import logging
import json
from typing import Dict

from pydantic import ValidationError

from hermes import hermes
from hermes.models import PrometheusHistogram
from hermes.hermes import push_udp_packet, set_hermes_config
from hermes.exceptions import InvalidMetricException

LOGGER = logging.getLogger('hermes.histograms')


def observe_histogram(metric_name: str, labels: Dict[str, str], observation: float):
    """Function used to push new UDP packet for
    Histogram instance defined on Hermes Server. Inputs
    are converted to Pydantic models for validation
    before they are converted to JSON and pushed
    to the Hermes Server

    Arguments:
        metric_name: str name of metric on Hermes Server
        labels: dict labels and their corresponding string
            values
        observation: float value to push to histogram
    """
    # set hermes configuration if not been set before
    if None in (hermes.HERMES_HOST, hermes.HERMES_PORT):
        LOGGER.warn('hermes configuration not set. setting host to default localhost:7789')
        set_hermes_config()
    try:
        LOGGER.debug('incrementing histogram metric \'%s\'', metric_name)
        payload = {'labels': labels, 'observation': observation}
        gauge = PrometheusHistogram(**{'metric_name': metric_name, 'payload': payload})
        push_udp_packet(json.loads(gauge.json()))
    except ValidationError:
        LOGGER.exception('received invalid histogram configuration')
        raise InvalidMetricException