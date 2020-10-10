"""Module demonstrating how the python client library can be used"""

from hermes import observe_summary


if __name__ == '__main__':

    labels = {'label_1': 'test-label'}
    observe_summary('sample_summary', labels=labels, observation=20.84)