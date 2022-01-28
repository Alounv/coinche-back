package domain

import (
	"time"
)

type ErrorMessage int

const (
	ErrGameFull ErrorMessage = iota
	ErrGameNotFound
)

type GameServiceType interface {
	ListGames() []Game
	GetGame(id int) Game
	CreateGame(name string) int
	JoinGame(id int, playerName string) error
}

type Game struct {
	Id        int
	Name      string
	Full      bool
	CreatedAt time.Time
	Players   []string
}
