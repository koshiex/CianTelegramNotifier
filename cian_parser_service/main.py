"""
Main entry point for the CIAN Parser microservice.
"""

import os
import signal
import sys

from .config import API_HOST, DEFAULT_PORT, LOG_LEVEL
from .cian_service import CianService
from .cache import CianCache
from .api import CianAPI
from .logger import setup_logging

def main():
    """Run the CIAN Parser microservice."""
    logger = setup_logging(LOG_LEVEL)
    
    cian_service = CianService()
    cian_cache = CianCache(cian_service)
    cian_api = CianAPI(cian_service, cian_cache)
    
    refresh_thread = cian_cache.start_refresh_thread()
    
    def signal_handler(sig, frame):
        logger.info("Shutting down...")
        cian_cache.stop_refresh_thread()
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    
    port = int(os.environ.get('PORT', DEFAULT_PORT))
    
    logger.info(f"Starting CIAN Parser microservice on {API_HOST}:{port}")
    cian_api.run(host=API_HOST, port=port)

if __name__ == "__main__":
    main()
