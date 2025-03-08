"""
Logging configuration for the CIAN Parser microservice.
"""

import logging
import sys
from logging.handlers import RotatingFileHandler
import os

os.makedirs('logs', exist_ok=True)

logger = logging.getLogger('cian_parser_service')

def setup_logging(log_level=logging.INFO):
    """
    Set up logging configuration.
    
    Args:
        log_level: The logging level to use.
    """
    logger.setLevel(log_level)
    logger.propagate = False
    
    if logger.handlers:
        logger.handlers.clear()
    
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(log_level)
    console_format = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    console_handler.setFormatter(console_format)
    
    file_handler = RotatingFileHandler(
        'logs/cian_parser_service.log',
        maxBytes=10485760,  # 10MB
        backupCount=5
    )
    file_handler.setLevel(log_level)
    file_format = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )
    file_handler.setFormatter(file_format)
    
    logger.addHandler(console_handler)
    logger.addHandler(file_handler)
    
    return logger
