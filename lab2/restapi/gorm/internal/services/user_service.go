package services

import (
	"errors"
	"gorm/internal/database"
	"gorm/internal/models"
)

func CreateUser(user models.User) error {
	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func UpdateUser(user models.User) error {
	result := database.DB.Model(&user).Updates(models.User{Name: user.Name, Age: user.Age})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func DeleteUser(id int) error {
	result := database.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
