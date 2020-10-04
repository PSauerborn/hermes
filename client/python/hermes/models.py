"""Module containing data modles for the Hermes library"""

from typing import Dict, Optional
from pydantic import BaseModel, validator, ValidationError


class GaugePayload(BaseModel):
    """Dataclass containing values for
    Gauge payload"""
    operation: str
    labels: Dict[str, str]
    value: Optional[float]

    @validator('operation')
    def check_operation(cls, val):
        """Function used to check that operation type
        is valid for gauges"""
        if val in ['increment', 'decrement', 'set']:
            return val
        raise ValidationError

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