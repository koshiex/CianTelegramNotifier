#!/bin/bash

# ЦИАН Telegram Notifier Stop Script

set -e

echo "🛑 Остановка ЦИАН Telegram Notifier..."

# Stop services
docker-compose down

echo "✅ Все сервисы остановлены!"
echo ""
echo "📊 Для полной очистки используйте:"
echo "   docker-compose down -v  # Удалить также volumes"
echo "   docker system prune     # Очистить неиспользуемые образы" 