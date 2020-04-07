package resolvers

import (
	"context"
	"log"

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
	panic("not implemented")
}
func (r *mutationResolver) UpdateEvent(ctx context.Context, input models.EventInput) (*models.Event, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteEvent(ctx context.Context, eventID string) (bool, error) {
	panic("not implemented")
}
