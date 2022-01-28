package domain

import (
	"errors"
	"time"
)

type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Players   []string
}

func (game Game) IsFull() bool {
	return len(game.Players) == 4
}

func (game *Game) AddPlayer(playerName string) error {
	var err error
	if len(game.Players) == 4 {
		err = errors.New("GAME IS FULL")
	}

	game.Players = append(game.Players, playerName)
	return err
}

func NewGame(name string, creatorName string) Game {
	return Game{
		Name:    name,
		Players: []string{creatorName},
	}
}
