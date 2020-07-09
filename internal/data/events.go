package data

import (
	"log"
	"time"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	gqlModel "github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"

	"github.com/StarWarsDev/legion-ops/internal/orm"
)

func FindEvents(orm *orm.ORM, max int, forUser *gqlModel.User) ([]event.Event, error) {
	db := orm.DB.New()
	var dbRecords []event.Event
	var count int
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Limit(max).
		Find(&dbRecords).
		Count(&count).
		Error
	if err != nil {
		log.Println(err)
		return dbRecords, err
	}

	return dbRecords, nil
}

func GetEventWithID(orm *orm.ORM, eventID string) (event.Event, error) {
	db := orm.DB.New()
	var dbEvent event.Event
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=?", eventID).
		First(&dbEvent).
		Error

	return dbEvent, err
}

func CreateEventWithInput(input *gqlModel.EventInput, organizer *user.User, orm *orm.ORM) (event.Event, error) {
	newEvent := event.Event{
		Name:      input.Name,
		Type:      input.Type.String(),
		Organizer: *organizer,
	}

	if input.HeadJudge != nil && *input.HeadJudge != "" {
		dbHeadJudge, err := GetUser(*input.HeadJudge, orm)
		if err != nil {
			return event.Event{}, err
		}
		newEvent.HeadJudge = &dbHeadJudge
	}

	for _, judgeID := range input.Judges {
		judge, err := GetUser(judgeID, orm)
		if err != nil {
			return newEvent, err
		}
		newEvent.Judges = append(newEvent.Judges, judge)
	}

	for _, playerID := range input.Players {
		player, err := GetUser(playerID, orm)
		if err != nil {
			return newEvent, err
		}
		newEvent.Players = append(newEvent.Players, player)
	}

	for _, day := range input.Days {
		start, err := time.Parse(time.RFC3339, day.StartAt)
		if err != nil {
			return newEvent, err
		}

		end, err := time.Parse(time.RFC3339, day.EndAt)
		if err != nil {
			return newEvent, err
		}

		eventDay := event.Day{
			StartAt: start.UTC().Unix(),
			EndAt:   end.UTC().Unix(),
		}

		for r, round := range day.Rounds {
			newRound := event.Round{
				Counter: r,
			}

			for _, match := range round.Matches {
				p1, err := GetUser(match.Player1, orm)
				if err != nil {
					return newEvent, err
				}

				p2, err := GetUser(match.Player1, orm)
				if err != nil {
					return newEvent, err
				}
				newMatch := event.Match{
					Player1: p1,
					Player2: p2,
				}

				if match.Player1MarginOfVictory != nil {
					newMatch.Player1MarginOfVictory = *match.Player1MarginOfVictory
				}

				if match.Player1VictoryPoints != nil {
					newMatch.Player1VictoryPoints = *match.Player1VictoryPoints
				}

				if match.Player2VictoryPoints != nil {
					newMatch.Player2VictoryPoints = *match.Player2VictoryPoints
				}

				if match.Player2MarginOfVictory != nil {
					newMatch.Player2MarginOfVictory = *match.Player2MarginOfVictory
				}

				if match.Blue != nil && *match.Blue != "" {
					blue, err := GetUser(*match.Blue, orm)
					if err != nil {
						return newEvent, err
					}
					newMatch.Blue = &blue
				}

				if match.Bye != nil && *match.Bye != "" {
					bye, err := GetUser(*match.Bye, orm)
					if err != nil {
						return newEvent, err
					}
					newMatch.Bye = &bye
				}

				if match.Winner != nil && *match.Winner != "" {
					winner, err := GetUser(*match.Winner, orm)
					if err != nil {
						return newEvent, err
					}
					newMatch.Winner = &winner
				}

				newRound.Matches = append(newRound.Matches, newMatch)
			}

			eventDay.Rounds = append(eventDay.Rounds, newRound)
		}

		newEvent.Days = append(newEvent.Days, eventDay)
	}

	db := orm.DB.New()
	err := db.Debug().Create(&newEvent).Error

	return newEvent, err
}
