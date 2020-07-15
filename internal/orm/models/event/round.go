package event

import (
	"time"

	"github.com/StarWarsDev/legion-ops/internal/constants"

	"github.com/StarWarsDev/legion-ops/internal/orm/models"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Round struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt int       `gorm:"not null"`
	UpdatedAt int       `gorm:"not null"`
	Counter   int       `gorm:"not null"`
	Day       Day       `gorm:"PRELOAD:false"`
	DayID     uuid.UUID
	Matches   []Match
}

func (record *Round) BeforeSave(scope *gorm.Scope) error {
	var err error
	if record.ID.String() == constants.BlankUUID {
		id, err := models.GenerateUUID()
		if err != nil {
			return err
		}

		err = scope.SetColumn("ID", id)
		if err != nil {
			return err
		}
	}

	unixNow := time.Now().UTC().Unix()

	if record.CreatedAt == 0 {
		err = scope.SetColumn("CreatedAt", unixNow)
		if err != nil {
			return err
		}
	}

	err = scope.SetColumn("UpdatedAt", unixNow)
	return err
}
