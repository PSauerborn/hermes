"""Module demonstrating how the python client library can be used"""

import time

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