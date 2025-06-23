package services

import (
	"telegram_bot_service/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateOrUpdateUser creates or updates a user
func (s *UserService) CreateOrUpdateUser(userID int64, username, firstName, lastName string) (*models.User, error) {
	user := &models.User{
		ID:        userID,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	result := s.db.Where("id = ?", userID).First(user)
	if result.Error == gorm.ErrRecordNotFound {
		// Create new user
		if err := s.db.Create(user).Error; err != nil {
			return nil, err
		}
	} else if result.Error != nil {
		return nil, result.Error
	} else {
		// Update existing user
		user.Username = username
		user.FirstName = firstName
		user.LastName = lastName
		user.IsActive = true
		if err := s.db.Save(user).Error; err != nil {
			return nil, err
		}
	}

	return user, nil
}

// GetUser gets a user by ID
func (s *UserService) GetUser(userID int64) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllActiveUsers gets all active users
func (s *UserService) GetAllActiveUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// DeactivateUser deactivates a user
func (s *UserService) DeactivateUser(userID int64) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error
}
