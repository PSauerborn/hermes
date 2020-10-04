"""Module demonstrating how the python client library can be used"""

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

    set_hermes_config('localost', 7789)
    example_job_wrapper_post()
