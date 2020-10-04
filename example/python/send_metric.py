"""Module used to demonstrate Hermes functionality"""

import socket
import json

SOCKET = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)

payload_counter = {
    'metric_name': 'sample_counter',
    'payload': {
        'labels': {
            'label_1': 'testing label 1'
        }
    }
}

payload_gauge = {
    'metric_name': 'sample_gauge',
    'payload': {
        'labels': {
            'label_1': 'testing label 1',
            'label_2': 'testing label 2'
        },
        'value': 67.4
    }
}

payload = json.dumps(payload_counter)
SOCKET.sendto(bytes(payload, encoding='utf-8'), ('localhost', 7789))

payload = json.dumps(payload_gauge)
SOCKET.sendto(bytes(payload, encoding='utf-8'), ('localhost', 7789))