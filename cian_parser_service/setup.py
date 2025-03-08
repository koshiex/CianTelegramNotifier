from setuptools import setup, find_packages

setup(
    name="cian_parser_service",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "flask>=2.0.1",
        "cianparser",
        "gunicorn>=20.1.0",
    ],
    entry_points={
        "console_scripts": [
            "cian-parser-service=cian_parser_service:main",
        ],
    },
)
