#!/bin/bash

# ЦИАН Telegram Notifier Logs Script

case "$1" in
    "bot"|"telegram")
        echo "📱 Логи Telegram Bot:"
        docker-compose logs -f telegram-bot
        ;;
    "parser"|"cian")
        echo "🔍 Логи CIAN Parser:"
        docker-compose logs -f cian-parser
        ;;
    "all"|"")
        echo "📊 Все логи:"
        docker-compose logs -f
        ;;
    *)
        echo "❌ Неизвестный параметр: $1"
        echo ""
        echo "📖 Использование:"
        echo "   ./scripts/logs.sh         # Все логи"
        echo "   ./scripts/logs.sh bot     # Логи бота"
        echo "   ./scripts/logs.sh parser  # Логи парсера"
        exit 1
        ;;
esac 