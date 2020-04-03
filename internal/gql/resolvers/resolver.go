package resolvers

import (
	"context"
	"log"
	"time"

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
	var records []*models.Event
	rctx := graphql.GetRequestContext(ctx)
	log.Println(rctx.RawQuery)
	db := r.ORM.DB.New()
	var dbRecords []models2.Event
	var count int
	err := db.Select("id, name, type, created_at, updated_at").Find(&dbRecords).Count(&count).Error
	log.Println(count)
	if err != nil {
		log.Println(err)
		return records, err
	}

	for _, dbEvent := range dbRecords {
		log.Printf(
			"ID: %s, CreatedAt: %d, UpdatedAt: %d, Name: %s, Type: %s\n",
			dbEvent.ID.String(),
			dbEvent.CreatedAt,
			dbEvent.UpdatedAt,
			dbEvent.Name,
			dbEvent.Type,
		)
		records = append(records, &models.Event{
			ID:        dbEvent.ID.String(),
			CreatedAt: time.Unix(dbEvent.CreatedAt, 0).String(),
			UpdatedAt: time.Unix(dbEvent.UpdatedAt, 0).String(),
			Name:      dbEvent.Name,
			Type:      models.EventType(dbEvent.Type),
		})
	}

	return records, nil
}
