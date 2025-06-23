package bot

import (
	"fmt"
	"strings"
	"telegram_bot_service/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// formatListingForDisplay formats a listing for display in Telegram
func (b *Bot) formatListingForDisplay(listing *models.Listing) string {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("🏠 *%s*\n", escapeMarkdown(listing.Title)))
	message.WriteString(fmt.Sprintf("💰 *%s*\n", listing.Price))

	if listing.Address != "" {
		message.WriteString(fmt.Sprintf("📍 %s\n", escapeMarkdown(listing.Address)))
	}

	if listing.Area != "" {
		message.WriteString(fmt.Sprintf("📐 Площадь: %s\n", listing.Area))
	}

	if listing.Rooms != "" {
		message.WriteString(fmt.Sprintf("🚪 Комнат: %s\n", listing.Rooms))
	}

	if listing.Floor != "" {
		message.WriteString(fmt.Sprintf("🏢 Этаж: %s\n", listing.Floor))
	}

	if listing.Metro != "" {
		message.WriteString(fmt.Sprintf("🚇 Метро: %s\n", listing.Metro))
	}

	if listing.Description != "" && len(listing.Description) > 0 {
		desc := listing.Description
		if len(desc) > 200 {
			desc = desc[:200] + "..."
		}
		message.WriteString(fmt.Sprintf("📝 %s\n", escapeMarkdown(desc)))
	}

	message.WriteString(fmt.Sprintf("🔗 [Смотреть на ЦИАН](%s)\n", listing.URL))

	if listing.PublishedAt != "" {
		message.WriteString(fmt.Sprintf("🕐 Опубликовано: %s", listing.PublishedAt))
	}

	return message.String()
}

// createListingsKeyboard creates inline keyboard for listings
func (b *Bot) createListingsKeyboard(listings []models.Listing, currentPage, totalPages int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Add favorite buttons for each listing
	for i, listing := range listings {
		favoriteButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("⭐ Добавить в избранное (%d)", i+1),
			fmt.Sprintf("fav_add:%s", listing.ID),
		)
		rows = append(rows, []tgbotapi.InlineKeyboardButton{favoriteButton})
	}

	// Add navigation buttons
	var navButtons []tgbotapi.InlineKeyboardButton

	if currentPage > 0 {
		navButtons = append(navButtons, tgbotapi.NewInlineKeyboardButtonData("⬅️ Предыдущая", fmt.Sprintf("listings_page:%d", currentPage-1)))
	}

	if currentPage < totalPages-1 {
		navButtons = append(navButtons, tgbotapi.NewInlineKeyboardButtonData("Следующая ➡️", fmt.Sprintf("listings_page:%d", currentPage+1)))
	}

	if len(navButtons) > 0 {
		rows = append(rows, navButtons)
	}

	// Add refresh button
	refreshButton := tgbotapi.NewInlineKeyboardButtonData("🔄 Обновить", "refresh_listings")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{refreshButton})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// createFavoritesKeyboard creates inline keyboard for favorites management
func (b *Bot) createFavoritesKeyboard(favorites []models.Favorite) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Add remove buttons for each favorite (max 10 to avoid too long keyboard)
	maxButtons := 10
	if len(favorites) > maxButtons {
		maxButtons = len(favorites)
	}

	for i := 0; i < maxButtons && i < len(favorites); i++ {
		favorite := favorites[i]
		removeButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🗑️ Удалить (%d)", i+1),
			fmt.Sprintf("fav_remove:%s", favorite.ListingID),
		)
		rows = append(rows, []tgbotapi.InlineKeyboardButton{removeButton})
	}

	// Add back to listings button
	backButton := tgbotapi.NewInlineKeyboardButtonData("📋 К объявлениям", "back_to_listings")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{backButton})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// escapeMarkdown escapes special markdown characters
func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}
