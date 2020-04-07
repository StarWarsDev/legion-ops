package mapper

import (
	"log"
	"time"

	"github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"
)

func GQLEvent(eventIn *event.Event) *models.Event {
	log.Printf(
		"ID: %s, CreatedAt: %d, UpdatedAt: %d, Name: %s, Type: %s\n",
		eventIn.ID.String(),
		eventIn.CreatedAt,
		eventIn.UpdatedAt,
		eventIn.Name,
		eventIn.Type,
	)
	eventOut := models.Event{
		ID:        eventIn.ID.String(),
		CreatedAt: time.Unix(eventIn.CreatedAt, 0).String(),
		UpdatedAt: time.Unix(eventIn.UpdatedAt, 0).String(),
		Name:      eventIn.Name,
		Type:      models.EventType(eventIn.Type),
		Organizer: GQLUser(&eventIn.Organizer),
	}

	if eventIn.HeadJudge != nil {
		eventOut.HeadJudge = GQLUser(eventIn.HeadJudge)
	}

	for _, judge := range eventIn.Judges {
		eventOut.Judges = append(eventOut.Judges, GQLUser(&judge))
	}

	for _, player := range eventIn.Players {
		eventOut.Players = append(eventOut.Players, GQLUser(&player))
	}

	for _, day := range eventIn.Days {
		dayOut := models.EventDay{
			CreatedAt: time.Unix(day.CreatedAt, 0).String(),
			EndAt:     time.Unix(day.EndAt, 0).String(),
			ID:        day.ID.String(),
			UpdatedAt: time.Unix(day.UpdatedAt, 0).String(),
			Rounds:    nil,
			StartAt:   time.Unix(day.StartAt, 0).String(),
		}

		for _, round := range day.Rounds {
			roundOut := models.Round{
				ID:      round.ID.String(),
				Counter: round.Counter,
				Matches: nil,
			}

			dayOut.Rounds = append(dayOut.Rounds, &roundOut)
		}

		eventOut.Days = append(eventOut.Days, &dayOut)
	}

	return &eventOut
}
