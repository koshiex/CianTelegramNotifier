"""
Configuration settings for the CIAN Parser microservice.
"""
import logging

# Default CIAN search settings
CIAN_SETTINGS = {
    "min_price": 30000,
    "max_price": 80000,
    "min_house_year": 1990,
    "max_house_year": 2023,
    "min_floor": 3,
    "sort_by": "total_meters_from_max_to_min",
}

# API settings
DEFAULT_PORT = 5000
API_HOST = '0.0.0.0'

# Cache settings
CACHE_REFRESH_INTERVAL = 1800  # 30 minutes in seconds

# Logging settings
LOG_LEVEL = logging.INFO
LOG_FORMAT = '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
LOG_FILE = 'logs/cian_parser_service.log'
LOG_MAX_BYTES = 10485760  # 10MB
LOG_BACKUP_COUNT = 5
