package jobs

import (
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

const firstUserUsername = "darthjeryyk"
const secondUserUsername = "masterjeryyk"

var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_USERS",
	Migrate: func(db *gorm.DB) error {
		firstUser := user.User{
			Username: firstUserUsername,
			Name:     "Darth Jeryyk",
			Picture:  "https://i.pinimg.com/originals/74/2a/93/742a935b67da21d21b46a44d2de4591a.jpg",
		}
		err := db.Debug().Create(&firstUser).Error
		if err != nil {
			return err
		}

		secondUser := user.User{
			Username: secondUserUsername,
			Name:     "Master Jeryyk",
			Picture:  "https://vignette.wikia.nocookie.net/starwars/images/6/61/Jedi_Master_Belth_Allusis.jpg/revision/latest?cb=20100220001352",
		}
		return db.Debug().Create(&secondUser).Error
	},
	Rollback: func(db *gorm.DB) error {
		var secondUser user.User
		db.First(&secondUser, &user.User{
			Username: secondUserUsername,
		})
		err := db.Debug().Delete(&secondUser).Error
		if err != nil {
			return err
		}

		var firstUser user.User
		db.First(&firstUser, &user.User{
			Username: firstUserUsername,
		})
		return db.Debug().Delete(&firstUser).Error
	},
}
