package event

import (
	"time"

	"github.com/StarWarsDev/legion-ops/internal/orm/models"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Day struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt int64     `gorm:"not null"`
	UpdatedAt int64     `gorm:"not null"`
	StartAt   int64     `gorm:"not null"`
	EndAt     int64     `gorm:"not null"`
	Event     Event     `gorm:"PRELOAD:false"`
	EventID   uuid.UUID
	Rounds    []Round
}

func (record *Day) BeforeCreate(scope *gorm.Scope) error {
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
