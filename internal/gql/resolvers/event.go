package resolvers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/StarWarsDev/legion-ops/routes/middlewares"

	"github.com/StarWarsDev/legion-ops/internal/gql/resolvers/mapper"

	"github.com/StarWarsDev/legion-ops/internal/data"
	"github.com/StarWarsDev/legion-ops/internal/gql/models"
)

// Query
func (r *queryResolver) Events(ctx context.Context, userID *string, max *int) ([]*models.Event, error) {
	var records []*models.Event

	// set some defaults and upper limits
	if max == nil {
		defaultMax := 10
		max = &defaultMax
	}

	if *max > 100 {
		hardMax := 100
		max = &hardMax
	}

	if userID != nil && *userID != "" {
		log.Println("Only getting records for the specified user")
	}

	dbRecords, err := data.FindEvents(r.ORM, *max, nil)
	if err != nil {
		return records, err
	}

	for _, dbEvent := range dbRecords {
		records = append(records, mapper.GQLEvent(&dbEvent))
	}

	return records, nil
}

// Mutation
func (r *mutationResolver) CreateEvent(ctx context.Context, input models.EventInput) (*models.Event, error) {
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot create event, valid user not supplied")
	}

	panic("not implemented")
}
func (r *mutationResolver) UpdateEvent(ctx context.Context, eventID string, input models.EventInput) (*models.Event, error) {
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot create event, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(r.ORM, eventID)
	if err != nil {
		return nil, err
	}

	if dbEvent.OrganizerID.String() != eventID {
		return nil, fmt.Errorf("account is not authorized to modify event")
	}

	panic("not implemented")
}
func (r *mutationResolver) DeleteEvent(ctx context.Context, eventID string) (bool, error) {
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return false, errors.New("cannot create event, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(r.ORM, eventID)
	if err != nil {
		return false, err
	}

	if dbEvent.OrganizerID.String() != eventID {
		return false, fmt.Errorf("account is not authorized to modify event")
	}

	panic("not implemented")
}
