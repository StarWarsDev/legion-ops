package resolvers

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql"

	"github.com/StarWarsDev/legion-ops/internal/orm"
	models2 "github.com/StarWarsDev/legion-ops/internal/orm/models"

	"github.com/StarWarsDev/legion-ops/internal/gql"
	"github.com/StarWarsDev/legion-ops/internal/gql/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	ORM *orm.ORM
}

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
	rctx := graphql.GetRequestContext(ctx)
	log.Println(rctx.RawQuery)
	db := r.ORM.DB.New()
	dbRecords := []models2.Event{}
	var count int
	db = db.Select("id, createdAt, lastUpdated, name, type").Find(&dbRecords).Count(&count)

	var records []*models.Event

	for _, dbEvent := range dbRecords {
		log.Println(count)
		log.Println(dbEvent)
		records = append(records, &models.Event{
			ID:          dbEvent.ID.String(),
			CreatedAt:   dbEvent.CreatedAt.String(),
			LastUpdated: dbEvent.LastUpdated.String(),
			Name:        dbEvent.Name,
			Type:        models.EventType(dbEvent.Type),
		})
	}

	return records, nil
}
