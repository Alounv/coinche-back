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
	Pause       Phase = 4
)

type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
	Players   []string
	Phase     Phase
}

const (
	ErrAlreadyInGame  = "ALREADY IN GAME"
	ErrGameFull       = "GAME IS FULL"
	ErrPlayerNotFound = "PLAYER NOT FOUND"
)

func (game Game) IsFull() bool {
	return len(game.Players) == 4
}

func (game *Game) AddPlayer(playerName string) error {
	for _, name := range game.Players {
		if name == playerName {
			return errors.New(ErrAlreadyInGame)
		}
	}

	if game.IsFull() {
		return errors.New(ErrGameFull)
	}

	game.Players = append(game.Players, playerName)
	if game.IsFull() && game.Phase == Preparation {
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
		return errors.New(ErrPlayerNotFound)
	}

	game.Players = newPlayers
	if !game.IsFull() && game.Phase != Preparation {
		game.Phase = Pause
	}
	return nil
}

func NewGame(name string, creatorName string) Game {
	return Game{
		Name:    name,
		Players: []string{creatorName},
		Phase:   Preparation,
	}
}
