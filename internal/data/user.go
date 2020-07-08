package data

import (
	"github.com/StarWarsDev/legion-ops/internal/orm"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"
)

func FindUserWithUsername(username string, orm *orm.ORM) (user.User, error) {
	db := orm.DB.New()
	var userRecord user.User
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("username=?", username).
		First(&userRecord).
		Error
	return userRecord, err
}
