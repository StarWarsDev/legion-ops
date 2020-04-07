package event

import (
	"time"

	"github.com/StarWarsDev/legion-ops/internal/orm/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Match struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt int       `gorm:"not null"`
	UpdatedAt int       `gorm:"not null"`

	// Match
	Round   Round
	RoundID uuid.UUID

	// Bye
	Bye   *user.User
	ByeID *uuid.UUID

	// Player 1
	Player1                user.User
	Player1ID              uuid.UUID
	Player1VictoryPoints   int
	Player1MarginOfVictory int

	// Player 2
	Player2                user.User
	Player2ID              uuid.UUID
	Player2VictoryPoints   int
	Player2MarginOfVictory int

	// Blue player
	Blue *user.User

	// Winner
	Winner   *user.User
	WinnerID *uuid.UUID
}

func (record *Match) BeforeCreate(scope *gorm.Scope) error {
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
