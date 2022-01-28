package domain

import (
	"time"
)

type GameServiceType interface {
	ListGames() []Game
	GetGame(id int) Game
	CreateGame(name string) int
	JoinGame(id int, playerName string) error
}

type Game struct {
	ID        int
	Name      string
	Full      bool
	CreatedAt time.Time
	Players   []string
}
