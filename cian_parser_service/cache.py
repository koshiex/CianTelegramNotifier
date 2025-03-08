"""
Cache management for CIAN data.
"""

import time
import threading
from .config import CACHE_REFRESH_INTERVAL
from .logger import logger

class CianCache:
    """Cache for CIAN listings data."""
    
    def __init__(self, cian_service):
        """
        Initialize the cache.
        
        Args:
            cian_service: Instance of CianService to fetch data.
        """
        self.cian_service = cian_service
        self.cache = {
            "listings": [],
            "last_updated": None
        }
        self._lock = threading.Lock()
        self._stop_event = threading.Event()
        logger.debug("CianCache initialized")
    
    def get_listings(self, force_refresh=False):
        """
        Get listings from cache or fetch new ones if needed.
        
        Args:
            force_refresh (bool): Force refresh of cache regardless of age.
            
        Returns:
            list: Cached listings.
        """
        with self._lock:
            if (force_refresh or 
                not self.cache["listings"] or 
                (self.cache["last_updated"] and 
                 time.time() - self.cache["last_updated"] > CACHE_REFRESH_INTERVAL)):
                try:
                    logger.info("Cache miss or force refresh, fetching new data")
                    self.cache["listings"] = self.cian_service.get_listings()
                    self.cache["last_updated"] = time.time()
                except Exception as e:
                    logger.error(f"Error refreshing cache: {e}")
                    if not self.cache["listings"]:
                        raise
            else:
                logger.debug("Returning cached listings")
            
            return self.cache["listings"]
    
    def invalidate(self):
        """Invalidate the cache, forcing a refresh on next access."""
        with self._lock:
            logger.info("Invalidating cache")
            self.cache["listings"] = []
            self.cache["last_updated"] = None
    
    def start_refresh_thread(self):
        """Start a background thread to periodically refresh the cache."""
        logger.info("Starting background refresh thread")
        self._stop_event.clear()
        refresh_thread = threading.Thread(
            target=self._refresh_periodically,
            daemon=True
        )
        refresh_thread.start()
        return refresh_thread
    
    def stop_refresh_thread(self):
        """Signal the refresh thread to stop."""
        logger.info("Stopping background refresh thread")
        self._stop_event.set()
    
    def _refresh_periodically(self):
        """Periodically refresh the cache in the background."""
        while not self._stop_event.is_set():
            try:
                logger.info("Refreshing cache in background...")
                listings = self.cian_service.get_listings()
                
                with self._lock:
                    self.cache["listings"] = listings
                    self.cache["last_updated"] = time.time()
                
                logger.info(f"Cache refreshed with {len(listings)} listings")
            except Exception as e:
                logger.error(f"Error in background refresh: {e}")
            
            logger.debug(f"Waiting {CACHE_REFRESH_INTERVAL} seconds until next refresh")
            self._stop_event.wait(CACHE_REFRESH_INTERVAL)
