package data

import (
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"
	"github.com/jinzhu/gorm"
)

func CreateUser(userIn user.User, db *gorm.DB) (user.User, error) {
	userOut := userIn
	err := db.Create(&userOut).Error
	return userOut, err
}

func FindUserWithUsername(username string, db *gorm.DB) (user.User, error) {
	var userRecord user.User
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("username=?", username).
		First(&userRecord).
		Error
	return userRecord, err
}

func GetUser(userID string, db *gorm.DB) (user.User, error) {
	var dbUser user.User
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=?", userID).
		First(&dbUser).
		Error

	return dbUser, err
}
