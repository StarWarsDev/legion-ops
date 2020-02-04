// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

type Event struct {
	CreatedAt   string      `json:"createdAt"`
	Days        []*EventDay `json:"days"`
	ID          string      `json:"id"`
	LastUpdated string      `json:"lastUpdated"`
	Name        string      `json:"name"`
}

type EventDay struct {
	CreatedAt   string   `json:"createdAt"`
	EndAt       string   `json:"endAt"`
	ID          string   `json:"id"`
	LastUpdated string   `json:"lastUpdated"`
	Rounds      []*Round `json:"rounds"`
	StartAt     string   `json:"startAt"`
}

type EventInput struct {
	Name string `json:"name"`
}

type Match struct {
	ID      string    `json:"id"`
	Players []*Player `json:"players"`
}

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	User *User  `json:"user"`
}

type Round struct {
	Counter int      `json:"counter"`
	ID      string   `json:"id"`
	Matches []*Match `json:"matches"`
}

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Picture  *string `json:"picture"`
	Username string  `json:"username"`
}