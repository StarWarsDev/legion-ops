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

	userID := dbUser.ID.String()
	events, err := r.Events(ctx, &userID, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if event.Organizer.ID == userID {
			profile.OrganizedEvents = append(profile.OrganizedEvents, event)
		}

		if event.HeadJudge != nil && event.HeadJudge.ID == userID {
			profile.JudgingEvents = append(profile.JudgingEvents, event)
		}

		for _, judge := range event.Judges {
			if judge != nil && judge.ID == userID {
				profile.JudgingEvents = append(profile.JudgingEvents, event)
			}
		}

		for _, player := range event.Players {
			if player != nil && player.ID == userID {
				profile.ParticipatingEvents = append(profile.ParticipatingEvents, event)
			}
		}
	}

	return &profile, nil
}
