package event

import (
	"html"
	"strings"
	"time"

	"github.com/StarWarsDev/legion-ops/internal/constants"

	"github.com/StarWarsDev/legion-ops/internal/orm/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID          uuid.UUID   `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string      `gorm:"not null"`
	Description string      `gorm:"not null;type:text;default:''"`
	Type        string      `gorm:"not null"`
	CreatedAt   int64       `gorm:"not null"`
	UpdatedAt   int64       `gorm:"not null"`
	Organizer   user.User   `gorm:"not null;association_autoupdate:false;association_autocreate:false"`
	OrganizerID uuid.UUID   `gorm:"not null"`
	Players     []user.User `gorm:"many2many:event_players;association_autoupdate:false;association_autocreate:false"`
	Judges      []user.User `gorm:"many2many:event_judges;association_autoupdate:false;association_autocreate:false"`
	HeadJudge   *user.User  `gorm:"association_autoupdate:false;association_autocreate:false"`
	HeadJudgeID *uuid.UUID
	Days        []Day
}

func (event *Event) BeforeSave(scope *gorm.Scope) error {
	var err error
	if event.ID.String() == constants.BlankUUID {
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

	if event.CreatedAt == 0 {
		err = scope.SetColumn("CreatedAt", unixNow)
		if err != nil {
			return err
		}
	}

	err = scope.SetColumn("UpdatedAt", unixNow)

	return err
}

func (event *Event) Prepare() {
	id, err := models.GenerateUUID()
	if err != nil {
		return
	}

	event.ID = id
	if event.CreatedAt == 0 {
		event.CreatedAt = time.Now().UTC().Unix()
	}
	event.UpdatedAt = time.Now().UTC().Unix()
	event.Name = html.EscapeString(strings.TrimSpace(event.Name))
	event.Type = html.EscapeString(strings.TrimSpace(event.Type))
}
