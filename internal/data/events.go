package data

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	gqlModel "github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"

	"github.com/StarWarsDev/legion-ops/internal/orm"
)

func FindEvents(db *gorm.DB, max int, forUser *gqlModel.User) ([]event.Event, error) {
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

func GetEventWithID(eventID string, db *gorm.DB) (event.Event, error) {
	var dbEvent event.Event
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=?", eventID).
		First(&dbEvent).
		Error

	return dbEvent, err
}

func GetDayWithIDForEvent(id string, evt *event.Event, db *gorm.DB) (event.Day, error) {
	var day event.Day
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=? and event_id=?", id, evt.ID.String()).
		First(&day).
		Error
	return day, err
}

func GetRoundWithIDForDay(id string, day *event.Day, db *gorm.DB) (event.Round, error) {
	var round event.Round
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=? and day_id=?", id, day.ID.String()).
		First(&round).
		Error
	return round, err
}

func GetMatchWithIDForRound(id string, round *event.Round, db *gorm.DB) (event.Match, error) {
	var match event.Match
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=? and round_id=?", id, round.ID.String()).
		First(&match).
		Error
	return match, err
}

func CreateEventWithInput(input *gqlModel.EventInput, organizer *user.User, orm *orm.ORM) (event.Event, error) {
	db := NewDB(orm)

	// the organizer can only be set during create
	dbEvent := event.Event{
		Organizer: *organizer,
		Name:      input.Name,
		Type:      input.Type.String(),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if input.HeadJudge != nil && *input.HeadJudge != "" {
			dbHeadJudge, err := GetUser(*input.HeadJudge, tx)
			if err != nil {
				return err
			}
			dbEvent.HeadJudge = &dbHeadJudge
		}

		for _, judgeID := range input.Judges {
			judge, err := GetUser(judgeID, tx)
			if err != nil {
				return err
			}
			dbEvent.Judges = append(dbEvent.Judges, judge)
		}

		for _, playerID := range input.Players {
			player, err := GetUser(playerID, tx)
			if err != nil {
				return err
			}
			dbEvent.Players = append(dbEvent.Players, player)
		}

		for _, day := range input.Days {
			start, err := time.Parse(time.RFC3339, day.StartAt)
			if err != nil {
				return err
			}

			end, err := time.Parse(time.RFC3339, day.EndAt)
			if err != nil {
				return err
			}

			eventDay := event.Day{
				StartAt: start.UTC().Unix(),
				EndAt:   end.UTC().Unix(),
			}

			for r, round := range day.Rounds {
				dbRound := event.Round{
					Counter: r,
				}

				for _, match := range round.Matches {
					p1, err := GetUser(match.Player1, tx)
					if err != nil {
						return err
					}

					p2, err := GetUser(match.Player2, tx)
					if err != nil {
						return err
					}
					dbMatch := event.Match{
						Player1: p1,
						Player2: p2,
					}

					if match.Player1MarginOfVictory != nil {
						dbMatch.Player1MarginOfVictory = *match.Player1MarginOfVictory
					}

					if match.Player1VictoryPoints != nil {
						dbMatch.Player1VictoryPoints = *match.Player1VictoryPoints
					}

					if match.Player2VictoryPoints != nil {
						dbMatch.Player2VictoryPoints = *match.Player2VictoryPoints
					}

					if match.Player2MarginOfVictory != nil {
						dbMatch.Player2MarginOfVictory = *match.Player2MarginOfVictory
					}

					if match.Blue != nil && *match.Blue != "" {
						blue, err := GetUser(*match.Blue, tx)
						if err != nil {
							return err
						}
						dbMatch.Blue = &blue
					}

					if match.Bye != nil && *match.Bye != "" {
						bye, err := GetUser(*match.Bye, tx)
						if err != nil {
							return err
						}
						dbMatch.Bye = &bye
					}

					if match.Winner != nil && *match.Winner != "" {
						winner, err := GetUser(*match.Winner, tx)
						if err != nil {
							return err
						}
						dbMatch.Winner = &winner
					}
					dbRound.Matches = append(dbRound.Matches, dbMatch)
				}
				eventDay.Rounds = append(eventDay.Rounds, dbRound)
			}
			dbEvent.Days = append(dbEvent.Days, eventDay)
		}

		// create the event
		err := tx.Create(&dbEvent).Error
		if err != nil {
			return err
		}

		// save the event
		err = tx.Save(&dbEvent).Error

		return err
	})

	return dbEvent, err
}

func UpdateEventWithInput(input *gqlModel.EventInput, orm *orm.ORM) (event.Event, error) {
	db := NewDB(orm)

	dbEvent, err := GetEventWithID(*input.ID, db)
	if err != nil {
		return dbEvent, err
	}

	// do everything inside a transaction to protect data in case of errors
	err = db.Transaction(func(tx *gorm.DB) error {
		if dbEvent.Name != input.Name {
			dbEvent.Name = input.Name
		}

		if input.HeadJudge != nil {
			if (dbEvent.HeadJudge == nil) || (dbEvent.HeadJudge != nil && dbEvent.HeadJudge.ID.String() != *input.HeadJudge) {
				headJudge, err := GetUser(*input.HeadJudge, tx)
				if err != nil {
					return err
				}
				dbEvent.HeadJudge = &headJudge
			}
		}

		err = tx.Save(&dbEvent).Error
		if err != nil {
			return err
		}

		// manage players
		var players []user.User
		for _, pID := range input.Players {
			player, err := GetUser(pID, tx)
			if err != nil {
				return err
			}
			players = append(players, player)
		}
		dbEvent.Players = players
		err = tx.Model(&dbEvent).Association("Players").Replace(players).Error
		if err != nil {
			return err
		}

		// manage judges
		var judges []user.User
		for _, jID := range input.Judges {
			judge, err := GetUser(jID, tx)
			if err != nil {
				return err
			}
			judges = append(judges, judge)
		}
		dbEvent.Judges = judges
		err = tx.Model(&dbEvent).Association("Judges").Replace(judges).Error
		if err != nil {
			return err
		}

		// manage days
		for _, dayIn := range input.Days {
			if dayIn == nil {
				return errors.New("day input is invalid")
			}

			if dayIn.ID == nil {
				return errors.New("day input id is missing")
			}
			day, err := GetDayWithIDForEvent(*dayIn.ID, &dbEvent, tx)
			if err != nil {
				return err
			}
			start, err := time.Parse(time.RFC3339, dayIn.StartAt)
			if err != nil {
				return err
			}

			end, err := time.Parse(time.RFC3339, dayIn.EndAt)
			if err != nil {
				return err
			}

			day.StartAt = start.UTC().Unix()
			day.EndAt = end.UTC().Unix()

			err = tx.Save(&day).Error
			if err != nil {
				return err
			}

			// manage rounds for this day
			for _, roundIn := range dayIn.Rounds {
				if roundIn == nil {
					return errors.New("round input is invalid")
				}

				if roundIn.ID == nil {
					return errors.New("round input id is missing")
				}

				round, err := GetRoundWithIDForDay(*roundIn.ID, &day, tx)
				if err != nil {
					return err
				}

				// TODO: allow for the counter to be modified here?

				// manage this round's matches
				for _, matchIn := range roundIn.Matches {
					if matchIn == nil {
						return errors.New("match input is invalid")
					}

					if matchIn.ID == nil {
						return errors.New("match input id is missing")
					}

					match, err := GetMatchWithIDForRound(*matchIn.ID, &round, tx)
					if err != nil {
						return err
					}

					player1, err := GetUser(matchIn.Player1, tx)
					if err != nil {
						return err
					}

					player2, err := GetUser(matchIn.Player2, tx)
					if err != nil {
						return err
					}

					// match players
					if match.Player1.ID != player1.ID {
						match.Player1 = player1
					}

					if match.Player2.ID != player2.ID {
						match.Player2 = player2
					}

					// player victory points
					if matchIn.Player1VictoryPoints != nil && match.Player1VictoryPoints != *matchIn.Player1VictoryPoints {
						match.Player1VictoryPoints = *matchIn.Player1VictoryPoints
					}

					if matchIn.Player2VictoryPoints != nil && match.Player2VictoryPoints != *matchIn.Player2VictoryPoints {
						match.Player2VictoryPoints = *matchIn.Player2VictoryPoints
					}

					// player margin of victory
					if matchIn.Player1MarginOfVictory != nil && match.Player1MarginOfVictory != *matchIn.Player1MarginOfVictory {
						match.Player1MarginOfVictory = *matchIn.Player1MarginOfVictory
					}

					if matchIn.Player2MarginOfVictory != nil && match.Player2MarginOfVictory != *matchIn.Player2MarginOfVictory {
						match.Player2MarginOfVictory = *matchIn.Player2MarginOfVictory
					}

					// blue player
					if matchIn.Blue != nil {
						blue, err := GetUser(*matchIn.Blue, tx)
						if err != nil {
							return err
						}

						if match.Blue == nil || match.Blue.ID != blue.ID {
							match.Blue = &blue
						}
					}

					// winner
					if matchIn.Winner != nil {
						winner, err := GetUser(*matchIn.Winner, tx)
						if err != nil {
							return err
						}

						if match.Winner == nil || match.Winner.ID != winner.ID {
							match.Winner = &winner
						}
					}

					// bye
					if matchIn.Bye != nil {
						bye, err := GetUser(*matchIn.Bye, tx)
						if err != nil {
							return err
						}

						if match.Bye == nil || match.Bye.ID != bye.ID {
							match.Bye = &bye
						}
					}

					err = tx.Save(&match).Error
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return dbEvent, err
	}

	dbEvent, err = GetEventWithID(*input.ID, db)

	return dbEvent, err
}

func DeleteEventWithID(id string, orm *orm.ORM) (bool, error) {
	db := NewDB(orm)

	err := db.Transaction(func(tx *gorm.DB) error {
		evt, err := GetEventWithID(id, tx)
		if err != nil {
			return err
		}

		err = tx.Delete(&evt).Error

		return err
	})

	return err == nil, err
}

// days

func CreateDay(dayInput *gqlModel.EventDayInput, eventID string, orm *orm.ORM) (event.Day, error) {
	var eventDay event.Day
	db := NewDB(orm)

	err := db.Transaction(func(tx *gorm.DB) error {
		dbEvent, err := GetEventWithID(eventID, tx)
		if err != nil {
			return err
		}

		start, err := time.Parse(time.RFC3339, dayInput.StartAt)
		if err != nil {
			return err
		}

		end, err := time.Parse(time.RFC3339, dayInput.EndAt)
		if err != nil {
			return err
		}

		eventDay = event.Day{
			StartAt: start.UTC().Unix(),
			EndAt:   end.UTC().Unix(),
			Event:   dbEvent,
		}

		for r, round := range dayInput.Rounds {
			dbRound := event.Round{
				Counter: r,
			}

			for _, match := range round.Matches {
				p1, err := GetUser(match.Player1, tx)
				if err != nil {
					return err
				}

				p2, err := GetUser(match.Player2, tx)
				if err != nil {
					return err
				}
				dbMatch := event.Match{
					Player1: p1,
					Player2: p2,
				}

				if match.Player1MarginOfVictory != nil {
					dbMatch.Player1MarginOfVictory = *match.Player1MarginOfVictory
				}

				if match.Player1VictoryPoints != nil {
					dbMatch.Player1VictoryPoints = *match.Player1VictoryPoints
				}

				if match.Player2VictoryPoints != nil {
					dbMatch.Player2VictoryPoints = *match.Player2VictoryPoints
				}

				if match.Player2MarginOfVictory != nil {
					dbMatch.Player2MarginOfVictory = *match.Player2MarginOfVictory
				}

				if match.Blue != nil && *match.Blue != "" {
					blue, err := GetUser(*match.Blue, tx)
					if err != nil {
						return err
					}
					dbMatch.Blue = &blue
				}

				if match.Bye != nil && *match.Bye != "" {
					bye, err := GetUser(*match.Bye, tx)
					if err != nil {
						return err
					}
					dbMatch.Bye = &bye
				}

				if match.Winner != nil && *match.Winner != "" {
					winner, err := GetUser(*match.Winner, tx)
					if err != nil {
						return err
					}
					dbMatch.Winner = &winner
				}
				dbRound.Matches = append(dbRound.Matches, dbMatch)
			}
			eventDay.Rounds = append(eventDay.Rounds, dbRound)
		}

		// save the day
		return tx.Create(&eventDay).Error
	})

	return eventDay, err
}
