package domain

import (
	"errors"
	"time"
)

type Phase int

const (
	Preparation Phase = 0
	Bidding     Phase = 1
	Playing     Phase = 2
	Counting    Phase = 3
)

type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Players   []string
	Phase     Phase
}

func (game Game) IsFull() bool {
	return len(game.Players) == 4
}

func (game *Game) AddPlayer(playerName string) error {
	if len(game.Players) == 4 {
		return errors.New("GAME IS FULL")
	}

	for _, name := range game.Players {
		if name == playerName {
			return errors.New("PLAYER NAME ALREADY IN GAME")
		}
	}

	game.Players = append(game.Players, playerName)
	if len(game.Players) == 4 {
		game.Phase = Bidding
	}
	return nil
}

func (game *Game) RemovePlayer(playerName string) error {
	newPlayers := []string{}
	for _, name := range game.Players {
		if name != playerName {
			newPlayers = append(newPlayers, name)
		}
	}

	if len(newPlayers) == len(game.Players) {
		return errors.New("PLAYER NOT FOUND")
	}

	game.Players = newPlayers
	return nil
}

func NewGame(name string, creatorName string) Game {
	return Game{
		Name:    name,
		Players: []string{creatorName},
		Phase:   Preparation,
	}
}
