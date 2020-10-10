import pathlib
from setuptools import setup, find_packages

setup(
    name="python-hermes",
    version="0.0.4",
    description="Python client library for pushing metrics to Hermes server instance",
    url="https://github.com/PSauerborn/hermes",
    author="Pascal Sauerborn",
    author_email="pascal.sauerborn@gmail.com",
    license="MIT",
    classifiers=[
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
    ],
    packages=find_packages(),
    include_package_data=True,
    install_requires=["pydantic"]
)