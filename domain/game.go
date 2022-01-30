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
	if len(game.Players) == 4 {
		return errors.New("GAME IS FULL")
	}

	for _, name := range game.Players { // the player is already in the game, just return
		if name == playerName {
			return nil
		}
	}

	game.Players = append(game.Players, playerName)
	return nil
}

func NewGame(name string, creatorName string) Game {
	return Game{
		Name:    name,
		Players: []string{creatorName},
	}
}
