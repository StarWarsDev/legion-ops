package jobs

import (
	"time"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	gqlModel "github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

const firstEventName = "Interstellar Upheaval"

var SeedEvents *gormigrate.Migration = &gormigrate.Migration{
	ID: "SEED_EVENTS",
	Migrate: func(db *gorm.DB) error {
		var firstUser user.User
		db.First(&firstUser, &user.User{
			Username: firstUserUsername,
		})

		var secondUser user.User
		db.First(&secondUser, &user.User{
			Username: secondUserUsername,
		})
		firstEvent := event.Event{
			Name:        firstEventName,
			Description: "Lets get ready to... Avoid line of sight and re-roll bad dice?",
			Type:        gqlModel.EventTypeFfgop.String(),
			Organizer:   firstUser,
			HeadJudge:   &firstUser,
			Players: []user.User{
				firstUser,
				secondUser,
			},
			Days: []event.Day{
				{
					StartAt: time.Now().UTC().Unix(),
					EndAt:   time.Now().Add(24 * time.Hour).UTC().Unix(),
					Rounds: []event.Round{
						{
							Counter: 1,
							Matches: []event.Match{
								{Player1: firstUser, Player2: secondUser},
							},
						},
					},
				},
			},
		}
		return db.Debug().Create(&firstEvent).Error
	},
	Rollback: func(db *gorm.DB) error {
		var firstEvent event.Event
		db.First(&firstEvent, &event.Event{
			Name: firstEventName,
		})
		return db.Debug().Delete(&firstEvent).Error
	},
}
