package resolvers

import (
	"context"

	"github.com/StarWarsDev/legion-ops/internal/gql"
	"github.com/StarWarsDev/legion-ops/internal/gql/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() gql.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() gql.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateEvent(ctx context.Context, input models.EventInput) (*models.Event, error) {
	panic("not implemented")
}
func (r *mutationResolver) UpdateEvent(ctx context.Context, input models.EventInput) (*models.Event, error) {
	panic("not implemented")
}
func (r *mutationResolver) DeleteEvent(ctx context.Context, eventID string) (bool, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Events(ctx context.Context) ([]*models.Event, error) {
	records := []*models.Event{
		&models.Event{
			ID:   "ec17af15-e354-440c-a09f-69715fc8b595",
			Name: "Mock Event",
		},
	}
	return records, nil
}
