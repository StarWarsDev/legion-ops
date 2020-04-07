package event

import (
	"html"
	"strings"
	"time"

	"github.com/StarWarsDev/legion-ops/internal/orm/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID          uuid.UUID   `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string      `gorm:"not null"`
	Type        string      `gorm:"not null"`
	CreatedAt   int64       `gorm:"not null"`
	UpdatedAt   int64       `gorm:"not null"`
	Organizer   user.User   `gorm:"not null"`
	OrganizerID uuid.UUID   `gorm:"not null"`
	Players     []user.User `gorm:"many2many:event_players"`
	Judges      []user.User `gorm:"many2many:event_judges"`
	HeadJudge   *user.User
	HeadJudgeID *uuid.UUID
	Days        []Day
}

func (event *Event) BeforeCreate(scope *gorm.Scope) error {
	id, err := models.GenerateUUID()
	if err != nil {
		return err
	}

	err = scope.SetColumn("ID", id)
	if err != nil {
		return err
	}
	unixNow := time.Now().UTC().Unix()
	err = scope.SetColumn("CreatedAt", unixNow)
	if err != nil {
		return err
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
	event.CreatedAt = time.Now().UTC().Unix()
	event.UpdatedAt = time.Now().UTC().Unix()
	event.Name = html.EscapeString(strings.TrimSpace(event.Name))
	event.Type = html.EscapeString(strings.TrimSpace(event.Type))
}
