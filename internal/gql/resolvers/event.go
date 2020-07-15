package resolvers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

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

	err := data.NewDB(r.ORM).Transaction(func(tx *gorm.DB) error {
		dbRecords, err := data.FindEvents(tx, *max, nil)
		if err != nil {
			return err
		}

		for _, dbEvent := range dbRecords {
			records = append(records, mapper.GQLEvent(&dbEvent))
		}

		return nil
	})

	return records, err
}

// Mutation
func (r *mutationResolver) CreateEvent(ctx context.Context, input models.EventInput) (*models.Event, error) {
	dbUser := middlewares.UserInContext(ctx)
	if dbUser == nil || dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot create event, valid user not supplied")
	}

	if input.Name == "" {
		return nil, errors.New("name cannot be blank")
	}

	dbEvent, err := data.CreateEventWithInput(&input, dbUser, r.ORM)
	if err != nil {
		return nil, err
	}

	eventOut := mapper.GQLEvent(&dbEvent)

	return eventOut, nil
}
func (r *mutationResolver) UpdateEvent(ctx context.Context, input models.EventInput) (*models.Event, error) {
	if input.ID == nil {
		return nil, errors.New("event id is required")
	}

	dbUser := middlewares.UserInContext(ctx)
	if dbUser == nil || dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot create event, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(*input.ID, data.NewDB(r.ORM))
	if err != nil {
		return nil, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		return nil, fmt.Errorf("account is not authorized to modify event")
	}

	dbEvent, err = data.UpdateEventWithInput(&input, r.ORM)
	if err != nil {
		return nil, err
	}

	eventOut := mapper.GQLEvent(&dbEvent)

	return eventOut, nil
}
func (r *mutationResolver) DeleteEvent(ctx context.Context, eventID string) (bool, error) {
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return false, errors.New("cannot create event, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return false, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		return false, fmt.Errorf("account is not authorized to modify event")
	}

	panic("not implemented")
}
