package resolvers

import (
	"context"
	"errors"

	"github.com/StarWarsDev/legion-ops/internal/gql/resolvers/mapper"
	"github.com/StarWarsDev/legion-ops/routes/middlewares"

	"github.com/StarWarsDev/legion-ops/internal/gql/models"
)

func (r *queryResolver) MyProfile(ctx context.Context) (*models.Profile, error) {
	dbUser := middlewares.UserInContext(ctx)
	if dbUser == nil || dbUser.Username == "" {
		// user cannot be nil and username cannot be blank, return an error
		return nil, errors.New("cannot get profile, valid user not supplied")
	}

	profile := models.Profile{
		Account: mapper.GQLUser(dbUser),
	}

	return &profile, nil
}
