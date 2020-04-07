package data

import (
	"log"

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
