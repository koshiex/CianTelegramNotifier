package services

import (
	"telegram_bot_service/internal/models"

	"gorm.io/gorm"
)

type FavoriteService struct {
	db *gorm.DB
}

func NewFavoriteService(db *gorm.DB) *FavoriteService {
	return &FavoriteService{db: db}
}

// AddToFavorites adds a listing to user's favorites
func (s *FavoriteService) AddToFavorites(userID int64, listing *models.Listing, note string) (*models.Favorite, error) {
	// Check if already in favorites
	var existing models.Favorite
	result := s.db.Where("user_id = ? AND listing_id = ?", userID, listing.ID).First(&existing)
	if result.Error == nil {
		// Already exists, update note if provided
		if note != "" {
			existing.Note = note
			if err := s.db.Save(&existing).Error; err != nil {
				return nil, err
			}
		}
		return &existing, nil
	} else if result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	// Create new favorite
	favorite := &models.Favorite{
		UserID:    userID,
		ListingID: listing.ID,
		Title:     listing.Title,
		Price:     listing.Price,
		URL:       listing.URL,
		Note:      note,
	}

	if err := s.db.Create(favorite).Error; err != nil {
		return nil, err
	}

	return favorite, nil
}

// RemoveFromFavorites removes a listing from user's favorites
func (s *FavoriteService) RemoveFromFavorites(userID int64, listingID string) error {
	return s.db.Where("user_id = ? AND listing_id = ?", userID, listingID).Delete(&models.Favorite{}).Error
}

// GetUserFavorites gets all favorites for a user
func (s *FavoriteService) GetUserFavorites(userID int64) ([]models.Favorite, error) {
	var favorites []models.Favorite
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

// GetFavorite gets a specific favorite
func (s *FavoriteService) GetFavorite(userID int64, listingID string) (*models.Favorite, error) {
	var favorite models.Favorite
	if err := s.db.Where("user_id = ? AND listing_id = ?", userID, listingID).First(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

// IsFavorite checks if a listing is in user's favorites
func (s *FavoriteService) IsFavorite(userID int64, listingID string) bool {
	var count int64
	s.db.Model(&models.Favorite{}).Where("user_id = ? AND listing_id = ?", userID, listingID).Count(&count)
	return count > 0
}

// UpdateFavoriteNote updates the note for a favorite
func (s *FavoriteService) UpdateFavoriteNote(userID int64, listingID, note string) error {
	return s.db.Model(&models.Favorite{}).Where("user_id = ? AND listing_id = ?", userID, listingID).Update("note", note).Error
}
