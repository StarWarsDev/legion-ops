package models

import (
	"html"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name      string    `gorm:"not null"`
	Type      string    `gorm:"not null"`
	CreatedAt int64     `gorm:"not null"`
	UpdatedAt int64     `gorm:"not null"`
}

func (event *Event) BeforeCreate(scope *gorm.Scope) error {
	id, err := GenerateUUID()
	if err != nil {
		return err
	}

	scope.SetColumn("ID", id)
	unixNow := time.Now().UTC().Unix()
	scope.SetColumn("CreatedAt", unixNow)
	scope.SetColumn("UpdatedAt", unixNow)
	return nil
}

func (event *Event) Prepare() {
	id, err := GenerateUUID()
	if err != nil {
		return
	}

	event.ID = id
	event.CreatedAt = time.Now().Unix()
	event.UpdatedAt = time.Now().Unix()
	event.Name = html.EscapeString(strings.TrimSpace(event.Name))
	event.Type = html.EscapeString(strings.TrimSpace(event.Type))
}

func GenerateUUID() (uuid.UUID, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, err
	}
	return id, nil
}
