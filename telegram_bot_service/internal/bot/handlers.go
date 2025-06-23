package bot

import (
	"fmt"
	"strconv"
	"strings"
	"telegram_bot_service/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (b *Bot) handleListingsCommand(chatID int64) {
	listings, err := b.cianService.GetListings(false)
	if err != nil {
		logrus.WithError(err).Error("Failed to get listings")
		b.sendMessage(chatID, "❌ Ошибка при получении объявлений. Попробуйте позже.")
		return
	}

	if len(listings) == 0 {
		b.sendMessage(chatID, "📭 Объявления не найдены.")
		return
	}

	// Send listings with pagination
	b.sendListingsPage(chatID, listings, 0)
}

func (b *Bot) sendListingsPage(chatID int64, listings []models.Listing, page int) {
	pageSize := 5
	totalPages := (len(listings) + pageSize - 1) / pageSize

	if page < 0 || page >= totalPages {
		return
	}

	start := page * pageSize
	end := start + pageSize
	if end > len(listings) {
		end = len(listings)
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("🏠 *Объявления (%d-%d из %d)*\n\n", start+1, end, len(listings)))

	for i := start; i < end; i++ {
		listing := listings[i]
		message.WriteString(b.formatListingForDisplay(&listing))
		message.WriteString("\n---\n\n")
	}

	// Create keyboard with navigation and favorite buttons
	keyboard := b.createListingsKeyboard(listings[start:end], page, totalPages)

	msg := tgbotapi.NewMessage(chatID, message.String())
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard

	if _, err := b.api.Send(msg); err != nil {
		logrus.WithError(err).Error("Failed to send listings")
	}
}

func (b *Bot) handleFavoritesCommand(chatID int64, userID int64) {
	favorites, err := b.favoriteService.GetUserFavorites(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get favorites")
		b.sendMessage(chatID, "❌ Ошибка при получении избранного.")
		return
	}

	if len(favorites) == 0 {
		b.sendMessage(chatID, "⭐ У вас пока нет избранных объявлений.")
		return
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("⭐ *Ваши избранные объявления (%d)*\n\n", len(favorites)))

	for i, favorite := range favorites {
		message.WriteString(fmt.Sprintf("*%d.* [%s](%s)\n", i+1, favorite.Title, favorite.URL))
		message.WriteString(fmt.Sprintf("💰 %s\n", favorite.Price))
		if favorite.Note != "" {
			message.WriteString(fmt.Sprintf("📝 *Заметка:* %s\n", favorite.Note))
		}
		message.WriteString(fmt.Sprintf("🕒 Добавлено: %s\n\n", favorite.CreatedAt.Format("02.01.2006 15:04")))
	}

	// Create keyboard for managing favorites
	keyboard := b.createFavoritesKeyboard(favorites)

	msg := tgbotapi.NewMessage(chatID, message.String())
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard
	msg.DisableWebPagePreview = true

	if _, err := b.api.Send(msg); err != nil {
		logrus.WithError(err).Error("Failed to send favorites")
	}
}

func (b *Bot) handleSettingsCommand(chatID int64) {
	settings, err := b.cianService.GetSettings()
	if err != nil {
		logrus.WithError(err).Error("Failed to get settings")
		b.sendMessage(chatID, "❌ Ошибка при получении настроек.")
		return
	}

	var message strings.Builder
	message.WriteString("⚙️ *Текущие настройки поиска:*\n\n")

	for key, value := range settings {
		message.WriteString(fmt.Sprintf("• *%s:* %v\n", key, value))
	}

	message.WriteString("\n💡 Для изменения настроек используйте команды или обратитесь к администратору.")

	b.sendMessage(chatID, message.String())
}

func (b *Bot) handleSubscribeCommand(chatID int64, userID int64) {
	// TODO: Implement subscription logic
	b.sendMessage(chatID, "🔔 Подписка на уведомления активирована! Вы будете получать уведомления о новых объявлениях.")
}

func (b *Bot) handleUnsubscribeCommand(chatID int64, userID int64) {
	// TODO: Implement unsubscription logic
	b.sendMessage(chatID, "🔕 Подписка на уведомления отключена.")
}

func (b *Bot) handleTextMessage(message *tgbotapi.Message) {
	// Handle non-command text messages
	chatID := message.Chat.ID
	b.sendMessage(chatID, "Используйте команды для взаимодействия с ботом. Напишите /help для справки.")
}

func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	userID := query.From.ID
	data := query.CallbackData

	// Acknowledge the callback query
	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		logrus.WithError(err).Error("Failed to acknowledge callback query")
	}

	parts := strings.Split(data, ":")
	if len(parts) < 2 {
		// Handle single action callbacks
		switch data {
		case "refresh_listings":
			b.handleRefreshListings(chatID)
		case "back_to_listings":
			b.handleListingsCommand(chatID)
		}
		return
	}

	action := parts[0]
	param := parts[1]

	switch action {
	case "fav_add":
		b.handleAddToFavorites(chatID, userID, param)
	case "fav_remove":
		b.handleRemoveFromFavorites(chatID, userID, param)
	case "listings_page":
		if page, err := strconv.Atoi(param); err == nil {
			// Get fresh listings and show page
			if listings, err := b.cianService.GetListings(false); err == nil {
				b.sendListingsPage(chatID, listings, page)
			}
		}
	}
}

func (b *Bot) handleAddToFavorites(chatID int64, userID int64, listingID string) {
	// Get listing details from CIAN API
	listings, err := b.cianService.GetListings(false)
	if err != nil {
		b.sendMessage(chatID, "❌ Ошибка при получении данных объявления.")
		return
	}

	var targetListing *models.Listing
	for _, listing := range listings {
		if listing.ID == listingID {
			targetListing = &listing
			break
		}
	}

	if targetListing == nil {
		b.sendMessage(chatID, "❌ Объявление не найдено.")
		return
	}

	_, err = b.favoriteService.AddToFavorites(userID, targetListing, "")
	if err != nil {
		logrus.WithError(err).Error("Failed to add to favorites")
		b.sendMessage(chatID, "❌ Ошибка при добавлении в избранное.")
		return
	}

	b.sendMessage(chatID, fmt.Sprintf("⭐ Объявление \"%s\" добавлено в избранное!", targetListing.Title))
}

func (b *Bot) handleRemoveFromFavorites(chatID int64, userID int64, listingID string) {
	err := b.favoriteService.RemoveFromFavorites(userID, listingID)
	if err != nil {
		logrus.WithError(err).Error("Failed to remove from favorites")
		b.sendMessage(chatID, "❌ Ошибка при удалении из избранного.")
		return
	}

	b.sendMessage(chatID, "🗑️ Объявление удалено из избранного.")
}

func (b *Bot) handleRefreshListings(chatID int64) {
	listings, err := b.cianService.GetListings(true) // Force refresh
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh listings")
		b.sendMessage(chatID, "❌ Ошибка при обновлении объявлений.")
		return
	}

	if len(listings) == 0 {
		b.sendMessage(chatID, "📭 Объявления не найдены.")
		return
	}

	b.sendMessage(chatID, "🔄 Объявления обновлены!")
	b.sendListingsPage(chatID, listings, 0)
}
