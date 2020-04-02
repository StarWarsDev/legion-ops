package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

type Event struct {
	//gorm.Model
	ID          uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time `gorm:"index;not null;default:CURRENT_TIMESTAMP"`
	LastUpdated time.Time `gorm:"index;not null;default:CURRENT_TIMESTAMP"`
	Name        string    `gorm:"not null"`
	Type        string    `gorm:"not null"`
}

func (event *Event) BeforeCreate(scope *gorm.Scope) error {
	id, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		return err
	}

	scope.SetColumn("ID", id)
	return nil
}
