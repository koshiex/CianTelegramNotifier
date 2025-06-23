# ЦИАН Telegram Notifier

Микросервисная система для получения уведомлений о новых объявлениях аренды недвижимости с сайта ЦИАН через Telegram бота.

## Архитектура

Система состоит из двух микросервисов:

1. **ЦИАН Parser Service** (Python) - парсит объявления с ЦИАН и предоставляет REST API
2. **Telegram Bot Service** (Go) - телеграм бот для взаимодействия с пользователями

## Функциональность

### Парсер ЦИАН (Python)
- 📊 Парсинг объявлений с ЦИАН
- 🔄 Кэширование результатов
- ⚙️ Настройка параметров поиска
- 🚀 REST API для интеграции

### Telegram Bot (Go)
- 📋 Просмотр актуальных объявлений
- ⭐ Добавление объявлений в избранное
- 🔔 Подписка на уведомления о новых предложениях
- ⚙️ Настройка параметров поиска через бота
- 📝 Добавление заметок к избранным объявлениям

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Telegram Bot Token (получить у [@BotFather](https://t.me/botfather))

### Настройка

1. **Клонируйте репозиторий:**
   ```bash
   git clone <repository-url>
   cd CianTelegramNotifier
   ```

2. **Настройте переменные окружения:**
   ```bash
   # Отредактируйте файл docker.env
   nano docker.env
   
   # Замените тестовый токен на реальный:
   TELEGRAM_TOKEN=your_real_telegram_bot_token_here
   ```

3. **Запустите систему (рекомендуемый способ):**
   ```bash
   ./scripts/start.sh
   ```

   Или вручную:
   ```bash
   docker-compose up --build -d
   ```

### Команды бота

- `/start` - Начать работу с ботом
- `/help` - Показать справку
- `/listings` - Показать текущие объявления
- `/favorites` - Показать избранные объявления
- `/settings` - Показать и настроить параметры поиска
- `/subscribe` - Подписаться на уведомления
- `/unsubscribe` - Отписаться от уведомлений

## API

### ЦИАН Parser Service

**Базовый URL:** `http://localhost:5000`

#### Endpoints

- `GET /listings` - Получить список объявлений
  - `?refresh=true` - принудительное обновление
- `GET /settings` - Получить текущие настройки поиска
- `PUT /settings` - Обновить настройки поиска
- `GET /health` - Проверка состояния сервиса

### Пример настроек поиска

```json
{
  "deal_type": "rent",
  "offer_type": "flat",
  "price_min": 30000,
  "price_max": 80000,
  "rooms": [1, 2],
  "metro": ["Сокольники", "Преображенская площадь"]
}
```

## Разработка

### Структура проекта

```
CianTelegramNotifier/
├── cian_parser_service/      # Python парсер
│   ├── api.py               # REST API
│   ├── cian_service.py      # Логика парсинга
│   ├── cache.py             # Кэширование
│   └── ...
├── telegram_bot_service/     # Go телеграм бот
│   ├── main.go              # Точка входа
│   ├── internal/
│   │   ├── bot/            # Логика бота
│   │   ├── services/       # Бизнес-логика
│   │   ├── models/         # Модели данных
│   │   └── config/         # Конфигурация
│   └── ...
└── docker-compose.yml       # Оркестрация
```

### Локальная разработка

#### Python сервис
```bash
cd cian_parser_service
pip install -r requirements.txt
python main.py
```

#### Go сервис
```bash
cd telegram_bot_service
go mod tidy
go run main.go
```

## Настройка в production

1. **Безопасность:**
   - Используйте крепкие пароли
   - Настройте HTTPS для API
   - Ограничьте доступ к API

2. **Мониторинг:**
   - Настройте логирование
   - Добавьте метрики
   - Настройте алерты

3. **Масштабирование:**
   - Используйте внешнюю базу данных
   - Настройте балансировку нагрузки
   - Кэшируйте данные

## Управление системой

### Удобные скрипты

```bash
# Запуск системы
./scripts/start.sh

# Остановка системы  
./scripts/stop.sh

# Просмотр логов
./scripts/logs.sh           # Все логи
./scripts/logs.sh bot       # Логи бота
./scripts/logs.sh parser    # Логи парсера

# Проверка статуса
./scripts/status.sh
```

### Health Check

Система включает в себя health check эндпоинты:

```bash
# Статус Telegram Bot
curl http://localhost:8080/health

# Статус CIAN Parser  
curl http://localhost:5000/health

# Готовность бота
curl http://localhost:8080/ready
```

## Troubleshooting

### Частые проблемы

1. **Бот не отвечает:**
   - Проверьте токен бота в `docker.env`
   - Убедитесь, что сервисы запущены: `./scripts/status.sh`
   - Проверьте логи: `./scripts/logs.sh bot`

2. **Не получаются объявления:**
   - Проверьте работу парсера: `./scripts/logs.sh parser`
   - Убедитесь, что ЦИАН доступен
   - Проверьте настройки поиска: `curl http://localhost:5000/settings`

3. **Ошибки базы данных:**
   - Проверьте права доступа к папке `./data`
   - Убедитесь, что SQLite доступен
   - Проверьте health check: `curl http://localhost:8080/health`

### Конфигурация

Основные переменные окружения в `docker.env`:

```env
# Обязательные
TELEGRAM_TOKEN=ваш_токен_бота

# Опциональные
LOG_LEVEL=info                    # debug, info, warn, error
CHECK_INTERVAL=10m               # Интервал проверки новых объявлений
HEALTH_CHECK_ENABLED=true        # Включить health check
HEALTH_CHECK_PORT=8080           # Порт для health check
```

## Лицензия

MIT License - см. файл [LICENSE](LICENSE)

## Вклад в проект

1. Форкните проект
2. Создайте ветку для новой функции (`git checkout -b feature/AmazingFeature`)
3. Зафиксируйте изменения (`git commit -m 'Add some AmazingFeature'`)
4. Отправьте ветку (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request 