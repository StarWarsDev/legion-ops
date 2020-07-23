package data

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	gqlModel "github.com/StarWarsDev/legion-ops/internal/gql/models"
	"github.com/StarWarsDev/legion-ops/internal/orm/models/event"

	"github.com/StarWarsDev/legion-ops/internal/orm"
)

func FindEvents(db *gorm.DB, max int, forUser *user.User, eventType *gqlModel.EventType, startsAfter, endsBefore *string) ([]event.Event, error) {
	var dbRecords []event.Event
	var count int

	var where []string
	var params []interface{}

	if eventType != nil {
		where = append(where, "type = ?")
		params = append(params, eventType.String())
	}

	if startsAfter != nil || endsBefore != nil {
		ids, err := eventIdsInRange(db, startsAfter, endsBefore)
		if err != nil {
			return dbRecords, err
		}

		if len(ids) > 0 {
			where = append(where, "id IN (?)")
			params = append(params, ids)
		}
	}

	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where(strings.Join(where, " AND "), params...).
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

func eventIdsInRange(db *gorm.DB, startsAfter *string, endsBefore *string) ([]string, error) {
	var ids []string

	if startsAfter != nil && *startsAfter != "" && endsBefore == nil {
		t, err := time.Parse(time.RFC3339, *startsAfter)
		if err != nil {
			return nil, err
		}

		var days []event.Day
		err = db.Select("DISTINCT event_id").Where("start_at >= ?", t.Unix()).Find(&days).Error
		if err != nil {
			return nil, err
		}
		for _, day := range days {
			ids = append(ids, day.EventID.String())
		}
	}

	if endsBefore != nil && *endsBefore != "" && startsAfter == nil {
		t, err := time.Parse(time.RFC3339, *endsBefore)
		if err != nil {
			return nil, err
		}

		var days []event.Day
		err = db.Select("DISTINCT event_id").Where("end_at <= ?", t.Unix()).Find(&days).Error
		if err != nil {
			return nil, err
		}
		for _, day := range days {
			ids = append(ids, day.EventID.String())
		}
	}

	if startsAfter != nil && *startsAfter != "" && endsBefore != nil && *endsBefore != "" {
		startT, err := time.Parse(time.RFC3339, *startsAfter)
		if err != nil {
			return nil, err
		}

		endT, err := time.Parse(time.RFC3339, *endsBefore)
		if err != nil {
			return nil, err
		}

		var days []event.Day
		err = db.Select("DISTINCT event_id").Where("start_at >= ? AND end_at <= ?", startT.UTC().Unix(), endT.UTC().Unix()).Find(&days).Error
		if err != nil {
			return nil, err
		}
		for _, day := range days {
			ids = append(ids, day.EventID.String())
		}
	}

	return ids, nil
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

func GetDayWithID(id string, db *gorm.DB) (event.Day, error) {
	var day event.Day
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=?", id).
		First(&day).
		Error
	return day, err
}

func GetRoundWithID(id string, db *gorm.DB) (event.Round, error) {
	var round event.Round
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=?", id).
		First(&round).
		Error
	return round, err
}

func GetMatchWithID(id string, db *gorm.DB) (event.Match, error) {
	var match event.Match
	err := db.
		Set("gorm:auto_preload", true).
		Select("*").
		Where("id=?", id).
		First(&match).
		Error
	return match, err
}

func CreateEventWithInput(input *gqlModel.EventInput, organizer *user.User, orm *orm.ORM) (event.Event, error) {
	db := NewDB(orm)

	// the organizer can only be set during create
	dbEvent := event.Event{
		Organizer:   *organizer,
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type.String(),
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

		if dbEvent.Description != input.Description {
			dbEvent.Description = input.Description
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
			_, err = UpdateDay(dayIn, tx)
			if err != nil {
				return err
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

func UpdateDay(input *gqlModel.EventDayInput, db *gorm.DB) (event.Day, error) {
	var dayOut event.Day

	if input == nil {
		return dayOut, errors.New("day input is invalid")
	}

	if input.ID == nil {
		return dayOut, errors.New("day input id is missing")
	}
	day, err := GetDayWithID(*input.ID, db)
	if err != nil {
		return dayOut, err
	}
	start, err := time.Parse(time.RFC3339, input.StartAt)
	if err != nil {
		return dayOut, err
	}

	end, err := time.Parse(time.RFC3339, input.EndAt)
	if err != nil {
		return dayOut, err
	}

	day.StartAt = start.UTC().Unix()
	day.EndAt = end.UTC().Unix()

	err = db.Save(&day).Error
	if err != nil {
		return dayOut, err
	}

	// manage rounds for this day
	for _, roundIn := range input.Rounds {
		// there really isn't anything to change for a round
		// so lets skip down to the matches

		// manage this round's matches
		for _, matchIn := range roundIn.Matches {
			_, err := UpdateMatch(matchIn, db)
			if err != nil {
				return dayOut, err
			}
		}
	}

	return GetDayWithID(*input.ID, db)
}

func DeleteDay(id string, orm *orm.ORM) (bool, error) {
	db := NewDB(orm)

	err := db.Transaction(func(tx *gorm.DB) error {
		day, err := GetDayWithID(id, tx)
		if err != nil {
			return err
		}

		err = tx.Delete(&day).Error

		return err
	})

	return err == nil, err
}

// rounds

func CreateRound(roundInput *gqlModel.RoundInput, dayID string, orm *orm.ORM) (event.Round, error) {
	var newRound event.Round
	db := NewDB(orm)

	err := db.Transaction(func(tx *gorm.DB) error {
		day, err := GetDayWithID(dayID, tx)
		if err != nil {
			return err
		}

		newRound = event.Round{
			Counter: len(day.Rounds),
			Day:     day,
		}

		for _, match := range roundInput.Matches {
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
			newRound.Matches = append(newRound.Matches, dbMatch)
		}

		// save the round
		return tx.Create(&newRound).Error
	})

	if err != nil {
		return newRound, err
	}

	return newRound, nil
}

func DeleteRound(id string, orm *orm.ORM) (bool, error) {
	db := NewDB(orm)

	err := db.Transaction(func(tx *gorm.DB) error {
		round, err := GetRoundWithID(id, tx)
		if err != nil {
			return err
		}

		err = tx.Delete(&round).Error

		return err
	})

	return err == nil, err
}

// matches

func CreateMatch(matchInput *gqlModel.MatchInput, roundID string, orm *orm.ORM) (event.Match, error) {
	var newMatch event.Match
	db := NewDB(orm)

	err := db.Transaction(func(tx *gorm.DB) error {
		round, err := GetRoundWithID(roundID, tx)
		if err != nil {
			return err
		}

		p1, err := GetUser(matchInput.Player1, tx)
		if err != nil {
			return err
		}

		p2, err := GetUser(matchInput.Player2, tx)
		if err != nil {
			return err
		}
		newMatch = event.Match{
			Player1: p1,
			Player2: p2,
			Round:   round,
		}

		if matchInput.Player1MarginOfVictory != nil {
			newMatch.Player1MarginOfVictory = *matchInput.Player1MarginOfVictory
		}

		if matchInput.Player1VictoryPoints != nil {
			newMatch.Player1VictoryPoints = *matchInput.Player1VictoryPoints
		}

		if matchInput.Player2VictoryPoints != nil {
			newMatch.Player2VictoryPoints = *matchInput.Player2VictoryPoints
		}

		if matchInput.Player2MarginOfVictory != nil {
			newMatch.Player2MarginOfVictory = *matchInput.Player2MarginOfVictory
		}

		if matchInput.Blue != nil && *matchInput.Blue != "" {
			blue, err := GetUser(*matchInput.Blue, tx)
			if err != nil {
				return err
			}
			newMatch.Blue = &blue
		}

		if matchInput.Bye != nil && *matchInput.Bye != "" {
			bye, err := GetUser(*matchInput.Bye, tx)
			if err != nil {
				return err
			}
			newMatch.Bye = &bye
		}

		if matchInput.Winner != nil && *matchInput.Winner != "" {
			winner, err := GetUser(*matchInput.Winner, tx)
			if err != nil {
				return err
			}
			newMatch.Winner = &winner
		}

		// save the round
		return tx.Create(&newMatch).Error
	})

	if err != nil {
		return newMatch, err
	}

	return newMatch, nil
}

func UpdateMatch(input *gqlModel.MatchInput, db *gorm.DB) (event.Match, error) {
	var matchOut event.Match

	if input == nil {
		return matchOut, errors.New("match input is invalid")
	}

	if input.ID == nil {
		return matchOut, errors.New("match input id is missing")
	}

	match, err := GetMatchWithID(*input.ID, db)
	if err != nil {
		return matchOut, err
	}

	player1, err := GetUser(input.Player1, db)
	if err != nil {
		return matchOut, err
	}

	player2, err := GetUser(input.Player2, db)
	if err != nil {
		return matchOut, err
	}

	// match players
	if match.Player1.ID != player1.ID {
		match.Player1 = player1
	}

	if match.Player2.ID != player2.ID {
		match.Player2 = player2
	}

	// player victory points
	if input.Player1VictoryPoints != nil && match.Player1VictoryPoints != *input.Player1VictoryPoints {
		match.Player1VictoryPoints = *input.Player1VictoryPoints
	}

	if input.Player2VictoryPoints != nil && match.Player2VictoryPoints != *input.Player2VictoryPoints {
		match.Player2VictoryPoints = *input.Player2VictoryPoints
	}

	// player margin of victory
	if input.Player1MarginOfVictory != nil && match.Player1MarginOfVictory != *input.Player1MarginOfVictory {
		match.Player1MarginOfVictory = *input.Player1MarginOfVictory
	}

	if input.Player2MarginOfVictory != nil && match.Player2MarginOfVictory != *input.Player2MarginOfVictory {
		match.Player2MarginOfVictory = *input.Player2MarginOfVictory
	}

	// blue player
	if input.Blue != nil {
		blue, err := GetUser(*input.Blue, db)
		if err != nil {
			return matchOut, err
		}

		if match.Blue == nil || match.Blue.ID != blue.ID {
			match.Blue = &blue
		}
	}

	// winner
	if input.Winner != nil {
		winner, err := GetUser(*input.Winner, db)
		if err != nil {
			return matchOut, err
		}

		if match.Winner == nil || match.Winner.ID != winner.ID {
			match.Winner = &winner
		}
	}

	// bye
	if input.Bye != nil {
		bye, err := GetUser(*input.Bye, db)
		if err != nil {
			return matchOut, err
		}

		if match.Bye == nil || match.Bye.ID != bye.ID {
			match.Bye = &bye
		}
	}

	err = db.Save(&match).Error
	if err != nil {
		return matchOut, err
	}

	return GetMatchWithID(*input.ID, db)
}

func DeleteMatch(id string, orm *orm.ORM) (bool, error) {
	db := NewDB(orm)

	err := db.Transaction(func(tx *gorm.DB) error {
		match, err := GetMatchWithID(id, tx)
		if err != nil {
			return err
		}

		err = tx.Delete(&match).Error

		return err
	})

	return err == nil, err
}
