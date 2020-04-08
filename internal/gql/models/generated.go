// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
)

type Record interface {
	IsRecord()
}

type Event struct {
	ID        string      `json:"id"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	Name      string      `json:"name"`
	Type      EventType   `json:"type"`
	Days      []*EventDay `json:"days"`
	Organizer *User       `json:"organizer"`
	HeadJudge *User       `json:"headJudge"`
	Judges    []*User     `json:"judges"`
	Players   []*User     `json:"players"`
}

func (Event) IsRecord() {}

type EventDay struct {
	CreatedAt string   `json:"createdAt"`
	EndAt     string   `json:"endAt"`
	ID        string   `json:"id"`
	UpdatedAt string   `json:"updatedAt"`
	Rounds    []*Round `json:"rounds"`
	StartAt   string   `json:"startAt"`
}

func (EventDay) IsRecord() {}

type EventDayInput struct {
	EndAt   string        `json:"endAt"`
	Rounds  []*RoundInput `json:"rounds"`
	StartAt string        `json:"startAt"`
}

type EventInput struct {
	Name      string           `json:"name"`
	Type      EventType        `json:"type"`
	Days      []*EventDayInput `json:"days"`
	Organizer string           `json:"organizer"`
	HeadJudge *string          `json:"headJudge"`
	Judges    []string         `json:"judges"`
	Players   []string         `json:"players"`
}

type Match struct {
	ID                     string `json:"id"`
	Player1                *User  `json:"player1"`
	Player1VictoryPoints   int    `json:"player1VictoryPoints"`
	Player1MarginOfVictory int    `json:"player1MarginOfVictory"`
	Player2                *User  `json:"player2"`
	Player2VictoryPoints   int    `json:"player2VictoryPoints"`
	Player2MarginOfVictory int    `json:"player2MarginOfVictory"`
	Bye                    *User  `json:"bye"`
	Blue                   *User  `json:"blue"`
	Winner                 *User  `json:"winner"`
}

type MatchInput struct {
	Player1                string  `json:"player1"`
	Player1VictoryPoints   *int    `json:"player1VictoryPoints"`
	Player1MarginOfVictory *int    `json:"player1MarginOfVictory"`
	Player2                string  `json:"player2"`
	Player2VictoryPoints   *int    `json:"player2VictoryPoints"`
	Player2MarginOfVictory *int    `json:"player2MarginOfVictory"`
	Bye                    *string `json:"bye"`
	Blue                   *string `json:"blue"`
	Winner                 *string `json:"winner"`
}

type Round struct {
	Counter int      `json:"counter"`
	ID      string   `json:"id"`
	Matches []*Match `json:"matches"`
}

type RoundInput struct {
	Matches []*MatchInput `json:"matches"`
}

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Picture  *string `json:"picture"`
	Username string  `json:"username"`
}

type EventType string

const (
	EventTypeLeague EventType = "LEAGUE"
	EventTypeFfgop  EventType = "FFGOP"
	EventTypeOther  EventType = "OTHER"
)

var AllEventType = []EventType{
	EventTypeLeague,
	EventTypeFfgop,
	EventTypeOther,
}

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeLeague, EventTypeFfgop, EventTypeOther:
		return true
	}
	return false
}

func (e EventType) String() string {
	return string(e)
}

func (e *EventType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventType", str)
	}
	return nil
}

func (e EventType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
