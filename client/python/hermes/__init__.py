# __init__.py

__version__ = '0.0.2'

from hermes.gauges import set_gauge, increment_gauge, decrement_gauge, \
    gauge_wrapper, hermes_gauge
from hermes.counters import increment_counter, counter_wrapper
from hermes.hermes import set_hermes_config