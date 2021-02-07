"""Module containing exceptions raised by the Heres module"""


class InvalidMetricException(Exception):
    """Exception raised when invalid metrics
    are sent to the Hermes client"""

class HermesConfigurationException(Exception):
    """Exception raised when Invalid
    Hermes configuration is found"""
