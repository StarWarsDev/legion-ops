package resolvers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	"github.com/jinzhu/gorm"

	"github.com/StarWarsDev/legion-ops/routes/middlewares"

	"github.com/StarWarsDev/legion-ops/internal/gql/resolvers/mapper"

	"github.com/StarWarsDev/legion-ops/internal/data"
	"github.com/StarWarsDev/legion-ops/internal/gql/models"
)

// Query

func (r *queryResolver) CanModifyEvent(ctx context.Context, id string) (bool, error) {
	dbUser := middlewares.UserInContext(ctx)
	if dbUser == nil || dbUser.Username == "" {
		// username cannot be blank, return an error
		return false, errors.New("cannot check event permission, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(id, data.NewDB(r.ORM))
	if err != nil {
		return false, err
	}

	return dbUser.ID == dbEvent.Organizer.ID, nil
}

func (r *queryResolver) Event(ctx context.Context, id string) (*models.Event, error) {
	var eventOut *models.Event
	err := data.NewDB(r.ORM).Transaction(func(tx *gorm.DB) error {
		dbEvent, err := data.GetEventWithID(id, tx)
		if err != nil {
			return err
		}

		eventOut = mapper.GQLEvent(&dbEvent)

		return nil
	})

	return eventOut, err
}

func (r *queryResolver) Events(ctx context.Context, userID *string, max *int, eventType *models.EventType, startsAfter, endsBefore *string) ([]*models.Event, error) {
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

	err := data.NewDB(r.ORM).Transaction(func(tx *gorm.DB) error {
		var forUser *user.User

		onlyPublishedEvents := true

		if userID != nil && *userID != "" {
			log.Println("Only getting records for the specified user")
			dbUser, err := data.GetUser(*userID, tx)
			if err != nil {
				return err
			}

			forUser = &dbUser
			onlyPublishedEvents = false
		}

		dbRecords, err := data.FindEvents(tx, *max, forUser, eventType, startsAfter, endsBefore, &onlyPublishedEvents)
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
// events
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
		return nil, errors.New("cannot update event, valid user not supplied")
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
		return false, errors.New("cannot delete event, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return false, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		return false, fmt.Errorf("account is not authorized to modify event")
	}

	return data.DeleteEventWithID(eventID, r.ORM)
}

func (r *mutationResolver) PublishEvent(ctx context.Context, eventID string) (*models.Event, error) {
	return nil, nil
}

func (r *mutationResolver) UnpublishEvent(ctx context.Context, eventID string) (*models.Event, error) {
	return nil, nil
}

// days
func (r *mutationResolver) CreateDay(ctx context.Context, dayInput models.EventDayInput, eventID string) (*models.EventDay, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot create day, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return nil, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		log.Println(dbEvent.Organizer.Username, dbUser.Username)
		return nil, fmt.Errorf("account is not authorized to modify event")
	}

	// create the new day
	newDay, err := data.CreateDay(&dayInput, eventID, r.ORM)
	if err != nil {
		return nil, err
	}

	return mapper.GQLEventDay(&newDay), nil
}

func (r *mutationResolver) UpdateDay(ctx context.Context, dayInput models.EventDayInput, eventID string) (*models.EventDay, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot update day, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return nil, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		log.Println(dbEvent.Organizer.Username, dbUser.Username)
		return nil, fmt.Errorf("account is not authorized to modify event")
	}

	var dbDay event.Day
	db := data.NewDB(r.ORM)
	err = db.Transaction(func(tx *gorm.DB) error {
		dbDay, err = data.UpdateDay(&dayInput, tx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return mapper.GQLEventDay(&dbDay), nil
}

func (r *mutationResolver) DeleteDay(ctx context.Context, dayID, eventID string) (bool, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return false, errors.New("cannot delete day, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return false, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		return false, fmt.Errorf("account is not authorized to modify event")
	}

	return data.DeleteDay(dayID, r.ORM)
}

// rounds
func (r *mutationResolver) CreateRound(ctx context.Context, roundInput models.RoundInput, dayID, eventID string) (*models.Round, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot create round, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return nil, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		log.Println(dbEvent.Organizer.Username, dbUser.Username)
		return nil, fmt.Errorf("account is not authorized to modify event")
	}

	newRound, err := data.CreateRound(&roundInput, dayID, r.ORM)
	if err != nil {
		return nil, err
	}

	return mapper.GQLRound(&newRound), nil
}

func (r *mutationResolver) DeleteRound(ctx context.Context, roundID, eventID string) (bool, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return false, errors.New("cannot delete round, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return false, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		return false, fmt.Errorf("account is not authorized to modify event")
	}

	return data.DeleteRound(roundID, r.ORM)
}

// matches
func (r *mutationResolver) CreateMatch(ctx context.Context, matchInput models.MatchInput, roundID, eventID string) (*models.Match, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot create match, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return nil, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		log.Println(dbEvent.Organizer.Username, dbUser.Username)
		return nil, fmt.Errorf("account is not authorized to modify event")
	}

	newMatch, err := data.CreateMatch(&matchInput, roundID, r.ORM)
	if err != nil {
		return nil, err
	}

	return mapper.GQLMatch(&newMatch), nil
}

func (r *mutationResolver) UpdateMatch(ctx context.Context, matchInput models.MatchInput, eventID string) (*models.Match, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return nil, errors.New("cannot update match, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return nil, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		log.Println(dbEvent.Organizer.Username, dbUser.Username)
		return nil, fmt.Errorf("account is not authorized to modify event")
	}

	var dbMatch event.Match
	db := data.NewDB(r.ORM)
	err = db.Transaction(func(tx *gorm.DB) error {
		dbMatch, err = data.UpdateMatch(&matchInput, tx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return mapper.GQLMatch(&dbMatch), nil
}

func (r *mutationResolver) DeleteMatch(ctx context.Context, matchID, eventID string) (bool, error) {
	// check authorization against event ownership
	dbUser := middlewares.UserInContext(ctx)
	if dbUser.Username == "" {
		// username cannot be blank, return an error
		return false, errors.New("cannot delete match, valid user not supplied")
	}

	dbEvent, err := data.GetEventWithID(eventID, data.NewDB(r.ORM))
	if err != nil {
		return false, err
	}

	if dbEvent.Organizer.ID != dbUser.ID {
		return false, fmt.Errorf("account is not authorized to modify event")
	}

	return data.DeleteMatch(matchID, r.ORM)
}
