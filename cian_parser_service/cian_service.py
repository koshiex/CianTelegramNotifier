"""
Service for fetching data from CIAN using cianparser.
"""

import cianparser
from .config import CIAN_SETTINGS
from .logger import logger

class CianService:
    """Service for interacting with CIAN."""
    
    def __init__(self, settings=None):
        """
        Initialize the CIAN service.
        
        Args:
            settings (dict, optional): Search settings for CIAN. 
                                      Defaults to settings from config.
        """
        self.settings = settings or CIAN_SETTINGS.copy()
        logger.debug(f"CianService initialized with settings: {self.settings}")
    
    def get_listings(self):
        """
        Get apartment listings from CIAN.
        
        Returns:
            list: List of apartment listings.
        """
        logger.info("Fetching data from CIAN...")
        moscow_parser = cianparser.CianParser(location="Москва")
        data = moscow_parser.get_flats(
            deal_type="rent_long", 
            rooms=("all"), 
            with_saving_csv=False, 
            additional_settings=self.settings
        )
        logger.info(f"Retrieved {len(data)} listings from CIAN")
        return data
    
    def update_settings(self, new_settings):
        """
        Update search settings.
        
        Args:
            new_settings (dict): New search settings to apply.
            
        Returns:
            dict: Updated settings.
        """
        logger.info(f"Updating settings with: {new_settings}")
        self.settings.update(new_settings)
        logger.debug(f"Settings updated to: {self.settings}")
        return self.settings
    
    def get_settings(self):
        """
        Get current search settings.
        
        Returns:
            dict: Current search settings.
        """
        logger.debug("Retrieving current settings")
        return self.settings
