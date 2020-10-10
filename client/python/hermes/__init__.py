# __init__.py

__version__ = '0.0.4'

from hermes.gauges import set_gauge, increment_gauge, decrement_gauge, \
    gauge_wrapper, hermes_gauge
from hermes.counters import increment_counter, counter_wrapper
from hermes.histograms import observe_histogram
from hermes.summaries import observe_summary
from hermes.hermes import set_hermes_config