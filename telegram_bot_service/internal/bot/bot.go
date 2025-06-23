package bot

import (
	"telegram_bot_service/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	api             *tgbotapi.BotAPI
	cianService     *services.CianService
	userService     *services.UserService
	favoriteService *services.FavoriteService
}

func New(token string, cianService *services.CianService, userService *services.UserService, favoriteService *services.FavoriteService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = false
	logrus.WithField("username", api.Self.UserName).Info("Authorized on account")

	return &Bot{
		api:             api,
		cianService:     cianService,
		userService:     userService,
		favoriteService: favoriteService,
	}, nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			go b.handleCallbackQuery(update.CallbackQuery)
		}
	}

	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Create or update user
	_, err := b.userService.CreateOrUpdateUser(
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to create/update user")
		return
	}

	if message.IsCommand() {
		b.handleCommand(message)
	} else {
		b.handleTextMessage(message)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := message.Command()

	logrus.WithFields(logrus.Fields{
		"user_id":  message.From.ID,
		"username": message.From.UserName,
		"command":  command,
	}).Info("Received command")

	switch command {
	case "start":
		b.handleStartCommand(chatID)
	case "help":
		b.handleHelpCommand(chatID)
	case "listings":
		b.handleListingsCommand(chatID)
	case "favorites":
		b.handleFavoritesCommand(chatID, message.From.ID)
	case "settings":
		b.handleSettingsCommand(chatID)
	case "subscribe":
		b.handleSubscribeCommand(chatID, message.From.ID)
	case "unsubscribe":
		b.handleUnsubscribeCommand(chatID, message.From.ID)
	default:
		b.sendMessage(chatID, "Неизвестная команда. Используйте /help для списка доступных команд.")
	}
}

func (b *Bot) handleStartCommand(chatID int64) {
	welcomeText := `🏠 Добро пожаловать в бот уведомлений о недвижимости ЦИАН!

Этот бот поможет вам:
• 📋 Просматривать актуальные объявления
• ⭐ Добавлять объявления в избранное
• 🔔 Получать уведомления о новых предложениях
• ⚙️ Настраивать параметры поиска

Используйте /help для получения списка команд.`

	b.sendMessage(chatID, welcomeText)
}

func (b *Bot) handleHelpCommand(chatID int64) {
	helpText := `📖 Доступные команды:

/start - Начать работу с ботом
/help - Показать эту справку
/listings - Показать текущие объявления
/favorites - Показать избранные объявления
/settings - Показать и настроить параметры поиска
/subscribe - Подписаться на уведомления
/unsubscribe - Отписаться от уведомлений

💡 Tip: Вы можете добавлять объявления в избранное прямо из списка!`

	b.sendMessage(chatID, helpText)
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := b.api.Send(msg); err != nil {
		logrus.WithError(err).Error("Failed to send message")
	}
}
