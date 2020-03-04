package jobs

import (
	"github.com/StarWarsDev/legion-ops/internal/orm/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var firstEvent *models.Event = &models.Event{
	Name: "Test Event",
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
