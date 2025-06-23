#!/bin/bash

# ЦИАН Telegram Notifier Start Script

set -e

echo "🚀 Запуск ЦИАН Telegram Notifier..."

# Check if docker.env exists
if [ ! -f "docker.env" ]; then
    echo "❌ Файл docker.env не найден!"
    echo "📝 Создаем файл docker.env из шаблона..."
    cp docker.env docker.env.example || {
        echo "❌ Не удалось создать файл docker.env"
        echo "💡 Создайте файл docker.env вручную на основе docker.env.example"
        exit 1
    }
    echo "✅ Файл docker.env создан!"
    echo "⚠️  ВНИМАНИЕ: Не забудьте указать ваш TELEGRAM_TOKEN в файле docker.env"
fi

# Check if TELEGRAM_TOKEN is set
if grep -q "1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghi" docker.env; then
    echo "⚠️  ВНИМАНИЕ: Используется тестовый токен!"
    echo "📱 Получите реальный токен у @BotFather и замените его в docker.env"
    echo "⏸️  Нажмите Enter для продолжения или Ctrl+C для отмены..."
    read
fi

# Create necessary directories
echo "📁 Создание необходимых директорий..."
mkdir -p data logs

# Build and start services
echo "🔨 Сборка и запуск сервисов..."
docker-compose up --build -d

# Wait for services to start
echo "⏳ Ожидание запуска сервисов..."
sleep 10

# Check services status
echo "🔍 Проверка статуса сервисов..."

# Check CIAN Parser
if curl -s http://localhost:5000/health > /dev/null; then
    echo "✅ CIAN Parser запущен"
else
    echo "❌ CIAN Parser не отвечает"
fi

# Check Telegram Bot
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Telegram Bot запущен"
else
    echo "❌ Telegram Bot не отвечает"
fi

echo ""
echo "🎉 Система запущена!"
echo ""
echo "📊 Полезные команды:"
echo "   docker-compose logs -f                 # Просмотр логов"
echo "   docker-compose logs -f telegram-bot    # Логи бота"
echo "   docker-compose logs -f cian-parser     # Логи парсера"
echo "   curl http://localhost:8080/health      # Статус бота"
echo "   curl http://localhost:5000/health      # Статус парсера"
echo "   docker-compose stop                    # Остановка"
echo "   docker-compose down                    # Полная остановка"
echo ""
echo "🤖 Ваш бот готов к работе!"
echo "💬 Найдите его в Telegram и отправьте команду /start" 