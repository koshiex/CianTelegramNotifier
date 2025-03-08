"""
REST API for the CIAN Parser microservice.
"""

from flask import Flask, jsonify, request
from .logger import logger

class CianAPI:
    """API for CIAN Parser microservice."""
    
    def __init__(self, cian_service, cian_cache):
        """
        Initialize the API.
        
        Args:
            cian_service: Instance of CianService.
            cian_cache: Instance of CianCache.
        """
        self.app = Flask(__name__)
        self.cian_service = cian_service
        self.cian_cache = cian_cache
        
        self._register_routes()
        logger.debug("CianAPI initialized and routes registered")
    
    def _register_routes(self):
        """Register API routes."""
        
        @self.app.route('/listings', methods=['GET'])
        def get_listings():
            """API endpoint to get listings."""
            force_refresh = request.args.get('refresh', '').lower() == 'true'
            logger.info(f"GET /listings request received (force_refresh={force_refresh})")
            
            try:
                listings = self.cian_cache.get_listings(force_refresh=force_refresh)
                logger.info(f"Returning {len(listings)} listings")
                return jsonify(listings)
            except Exception as e:
                logger.error(f"Error getting listings: {e}", exc_info=True)
                return jsonify({"error": str(e)}), 500
        
        @self.app.route('/settings', methods=['GET'])
        def get_settings():
            """API endpoint to get current search settings."""
            logger.info("GET /settings request received")
            settings = self.cian_service.get_settings()
            logger.debug(f"Returning settings: {settings}")
            return jsonify(settings)
        
        @self.app.route('/settings', methods=['PUT'])
        def update_settings():
            """API endpoint to update search settings."""
            logger.info("PUT /settings request received")
            new_settings = request.json
            
            # Validate settings
            if not isinstance(new_settings, dict):
                logger.warning(f"Invalid settings format received: {type(new_settings)}")
                return jsonify({"error": "Settings must be a JSON object"}), 400
            
            logger.debug(f"Updating settings with: {new_settings}")
            # Update settings
            updated_settings = self.cian_service.update_settings(new_settings)
            
            # Invalidate cache
            self.cian_cache.invalidate()
            
            logger.info("Settings updated successfully")
            return jsonify({
                "message": "Settings updated", 
                "settings": updated_settings
            })
        
        @self.app.route('/health', methods=['GET'])
        def health_check():
            """Health check endpoint."""
            logger.debug("Health check request received")
            return jsonify({"status": "healthy"})
    
    def run(self, host, port, debug=False):
        """
        Run the Flask application.
        
        Args:
            host (str): Host to bind to.
            port (int): Port to bind to.
            debug (bool): Whether to run in debug mode.
        """
        logger.info(f"Starting Flask application on {host}:{port} (debug={debug})")
        self.app.run(host=host, port=port, debug=debug)
