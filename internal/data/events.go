package data

import (
	"log"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"

	gqlModel "github.com/StarWarsDev/legion-ops/internal/gql/models"

	"github.com/StarWarsDev/legion-ops/internal/orm"
)

func FindEvents(orm *orm.ORM, max int, forUser *gqlModel.User) ([]event.Event, error) {
	db := orm.DB.New()
	var dbRecords []event.Event
	var count int
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Limit(max).
		Find(&dbRecords).
		Count(&count).
		Error
	if err != nil {
		log.Println(err)
		return dbRecords, err
	}

	return dbRecords, nil
}

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
