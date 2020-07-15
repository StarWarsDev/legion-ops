package user

import (
	"html"
	"strings"
	"time"

	"github.com/StarWarsDev/legion-ops/internal/constants"

	"github.com/StarWarsDev/legion-ops/internal/orm/models"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// User represents a stored User record
type User struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username  string    `gorm:"not null;UNIQUE"`
	CreatedAt int64     `gorm:"not null"`
	UpdatedAt int64     `gorm:"not null"`
	Name      string    // if this is blank the UI will fall back to `Username`
	Picture   string
}

func (user *User) BeforeSave(scope *gorm.Scope) error {
	var err error
	if user.ID.String() == constants.BlankUUID {
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

	if user.CreatedAt == 0 {
		err = scope.SetColumn("CreatedAt", unixNow)
		if err != nil {
			return err
		}
	}

	err = scope.SetColumn("UpdatedAt", unixNow)
	return err
}

func (user *User) Prepare() {
	id, err := models.GenerateUUID()
	if err != nil {
		return
	}

	user.ID = id
	user.CreatedAt = time.Now().UTC().Unix()
	user.UpdatedAt = time.Now().UTC().Unix()
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Picture = strings.TrimSpace(user.Picture)
}

func (user *User) DisplayName() string {
	if user.Name != "" {
		return user.Name
	}

	return user.Username
}
