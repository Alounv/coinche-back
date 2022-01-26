package app

import "time"

type GameServiceType interface {
	ListGames() []Game
	GetGame(id int) Game
	CreateGame(name string) int
}

type Game struct {
	Id        int
	Name      string
	CreatedAt time.Time
}
