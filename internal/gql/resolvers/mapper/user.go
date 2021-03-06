package mapper

import (
	"github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"
)

func GQLUser(userIn *user.User) *models.User {
	if userIn == nil {
		return nil
	}
	picture := userIn.Picture
	return &models.User{
		ID:      userIn.ID.String(),
		Name:    userIn.DisplayName(),
		Picture: &picture,
	}
}
