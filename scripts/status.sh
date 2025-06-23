#!/bin/bash

# ЦИАН Telegram Notifier Status Script

echo "🔍 Проверка статуса ЦИАН Telegram Notifier..."
echo ""

# Check if containers are running
echo "📦 Статус контейнеров:"
docker-compose ps

echo ""
echo "🏥 Health Check:"

# Check CIAN Parser
echo -n "CIAN Parser: "
if curl -s http://localhost:5000/health > /dev/null 2>&1; then
    echo "✅ Работает"
    echo "   📊 Детали: $(curl -s http://localhost:5000/health | jq -r '.status // "OK"')"
else
    echo "❌ Не отвечает"
fi

# Check Telegram Bot  
echo -n "Telegram Bot: "
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Работает"
    echo "   📊 Детали: $(curl -s http://localhost:8080/health | jq -r '.status // "OK"')"
else
    echo "❌ Не отвечает"
fi

echo ""
echo "📈 Использование ресурсов:"
docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" cian-parser telegram-bot

echo ""
echo "📁 Размер данных:"
du -sh data/ logs/ 2>/dev/null || echo "Директории data/logs не найдены" 