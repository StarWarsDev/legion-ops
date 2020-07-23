package mapper

import (
	"time"

	"github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"
)

func GQLEvent(eventIn *event.Event) *models.Event {
	eventOut := models.Event{
		ID:          eventIn.ID.String(),
		CreatedAt:   time.Unix(eventIn.CreatedAt, 0).UTC().Format(time.RFC3339),
		UpdatedAt:   time.Unix(eventIn.UpdatedAt, 0).UTC().Format(time.RFC3339),
		Name:        eventIn.Name,
		Description: eventIn.Description,
		Type:        models.EventType(eventIn.Type),
		Organizer:   GQLUser(&eventIn.Organizer),
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
		eventOut.Days = append(eventOut.Days, GQLEventDay(&day))
	}

	return &eventOut
}

func GQLEventDay(day *event.Day) *models.EventDay {
	dayOut := models.EventDay{
		CreatedAt: time.Unix(day.CreatedAt, 0).UTC().Format(time.RFC3339),
		EndAt:     time.Unix(day.EndAt, 0).UTC().Format(time.RFC3339),
		ID:        day.ID.String(),
		UpdatedAt: time.Unix(day.UpdatedAt, 0).UTC().Format(time.RFC3339),
		Rounds:    nil,
		StartAt:   time.Unix(day.StartAt, 0).UTC().Format(time.RFC3339),
	}

	for _, round := range day.Rounds {
		dayOut.Rounds = append(dayOut.Rounds, GQLRound(&round))
	}

	return &dayOut
}

func GQLRound(round *event.Round) *models.Round {
	roundOut := models.Round{
		ID:      round.ID.String(),
		Counter: round.Counter,
		Matches: nil,
	}

	for _, match := range round.Matches {
		roundOut.Matches = append(roundOut.Matches, GQLMatch(&match))
	}

	return &roundOut
}

func GQLMatch(match *event.Match) *models.Match {
	return &models.Match{
		ID:                     match.ID.String(),
		Player1:                GQLUser(&match.Player1),
		Player1VictoryPoints:   match.Player1VictoryPoints,
		Player1MarginOfVictory: match.Player1MarginOfVictory,
		Player2:                GQLUser(&match.Player2),
		Player2VictoryPoints:   match.Player2VictoryPoints,
		Player2MarginOfVictory: match.Player2MarginOfVictory,
		Bye:                    GQLUser(match.Bye),
		Blue:                   GQLUser(match.Blue),
		Winner:                 GQLUser(match.Winner),
	}
}
