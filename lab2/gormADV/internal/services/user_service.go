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

func GetUsersWithProfiles(minAge, maxAge, page, pageSize int, sort string) ([]models.User, int, error) {
	var users []models.User
	var totalCount int64

	db := database.DB.Model(&models.User{})

	if minAge > 0 {
		db = db.Where("age >= ?", minAge)
	}
	if maxAge > 0 {
		db = db.Where("age <= ?", maxAge)
	}

	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, err
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
	if err := db.Preload("Profile").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(totalCount), nil
}
func UpdateUserAndProfile(user *models.User, profile *models.Profile) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {

		result := tx.Model(&models.User{}).
			Where("id = ?", user.ID).
			Updates(map[string]interface{}{
				"name": user.Name,
				"age":  user.Age,
			})
		if result.Error != nil {
			return fmt.Errorf("failed to update user: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("user not found")
		}

		if profile != nil {
			result = tx.Model(&models.Profile{}).
				Where("user_id = ?", user.ID).
				Updates(map[string]interface{}{
					"bio":                 profile.Bio,
					"profile_picture_url": profile.ProfilePictureURL,
				})
			if result.Error != nil {
				return fmt.Errorf("failed to update profile: %w", result.Error)
			}

		}
		return nil
	})
}

func DeleteUserWithProfile(userID uint) error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.User{}, userID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("user not found")
		}

		result = tx.Delete(&models.Profile{}, "user_id = ?", userID)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
