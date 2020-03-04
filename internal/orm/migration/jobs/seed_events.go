package jobs

import (
	m2 "github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var firstEvent *models.Event = &models.Event{
	Name: "Test Event",
	Type: m2.EventTypeLeague.String(),
}

var SeedEvents *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_EVENTS",
	Migrate: func(db *gorm.DB) error {
		return db.Create(firstEvent).Error
	},
	Rollback: func(db *gorm.DB) error {
		return db.Delete(firstEvent).Error
	},
}
