"""Module demonstrating how the python client library can be used"""

from hermes import observe_histogram


if __name__ == '__main__':

    labels = {'label_1': 'test-label'}
    observe_histogram('sample_histogram', labels=labels, observation=6.54)