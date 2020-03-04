package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `gorm:"index;not null;default:CURRENT_TIMESTAMP"`
	//UpdatedAt *time.Time `gorm:"index"`
}

type BaseModelSoftDelete struct {
	BaseModel
	//DeletedAt *time.Time `sql:"index"`
}
