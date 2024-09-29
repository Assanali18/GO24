package services

import (
	"fmt"
	"gorm.io/gorm"
	"gormADV/internal/database"
	"gormADV/internal/models"
	"log"
)

func CreateUserWithProfile(user *models.User) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to create user and profile: %v", err)
		return err
	}

	log.Println("User and profile created successfully.")
	return nil
}

func GetUsersWithProfiles(minAge, maxAge, page, pageSize int, sort string) ([]models.User, error) {
	var users []models.User
	db := database.DB

	if minAge > 0 {
		db = db.Where("age >= ?", minAge)
	}
	if maxAge > 0 {
		db = db.Where("age <= ?", maxAge)
	}

	switch sort {
	case "name_asc":
		db = db.Order("name ASC")
	case "name_desc":
		db = db.Order("name DESC")
	default:
		db = db.Order("id")
	}

	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)

	err := db.Preload("Profile").Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}
func UpdateUserAndProfile(user *models.User, profile *models.Profile) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.User{}).Where("id = ?", user.ID).Updates(user)
		if result.Error != nil {
			return fmt.Errorf("failed to update user: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("user not found")
		}

		result = tx.Model(&models.Profile{}).Where("user_id = ?", user.ID).Updates(profile)
		if result.Error != nil {
			return fmt.Errorf("failed to update profile: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("profile not found")
		}

		return nil
	})
}
func DeleteUserWithProfile(userID uint) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, userID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.Profile{}, "user_id = ?", userID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
