"""Module containing data modles for the Hermes library"""

from typing import Dict
from pydantic import BaseModel


class GaugePayload(BaseModel):
    """Dataclass containing values for
    Gauge payload"""
    value: float
    labels: Dict[str, str]

class PrometheusGauge(BaseModel):
    """Dataclass containing values for
    a Prometheus Gauge instance"""
    metric_name: str
    payload: GaugePayload

class CounterPayload(BaseModel):
    """Dataclass containing values for
    a Prometheus Counter instance"""
    labels: Dict[str, str]

class PrometheusCounter(BaseModel):
    """Dataclass containing values for
    a Prometheus counter instance"""
    metric_name: str
    payload: CounterPayload