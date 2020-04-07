package jobs

import (
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

const firstUserUsername = "darthjeryyk"

var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		firstUser := user.User{
			Username: firstUserUsername,
			Name:     "Darth Jeryyk",
			Picture:  "https://i.pinimg.com/originals/74/2a/93/742a935b67da21d21b46a44d2de4591a.jpg",
		}
		return db.Debug().Create(&firstUser).Error
	},
	Rollback: func(db *gorm.DB) error {
		var firstUser user.User
		db.First(&firstUser, &user.User{
			Username: firstUserUsername,
		})
		return db.Debug().Delete(&firstUser).Error
	},
}
