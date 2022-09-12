package domain

import (
	"errors"
)

const (
	ErrAlreadyInGame   = "ALREADY IN GAME"
	ErrEmptyPlayerName = "EMPTY PLAYER NAME"
	ErrGameFull        = "GAME IS FULL"
	ErrPlayerNotFound  = "PLAYER NOT FOUND"
)

func (game Game) IsFull() bool {
	return len(game.Players) == 4
}

func (game *Game) AddPlayer(playerName string) error {
	if playerName == "" {
		return errors.New(ErrEmptyPlayerName)
	}

	if _, ok := game.Players[playerName]; ok {
		return errors.New(ErrAlreadyInGame)
	}

	if game.IsFull() {
		return errors.New(ErrGameFull)
	}

	game.Players[playerName] = Player{}

	if game.IsFull() && (game.Phase == Preparation || game.Phase == Pause) {
		game.Phase = Teaming
	}
	return nil
}

func (game *Game) RemovePlayer(playerName string) error {
	if _, ok := game.Players[playerName]; !ok {
		return errors.New(ErrPlayerNotFound)
	}

	delete(game.Players, playerName)

	if !game.IsFull() && game.Phase != Preparation {
		game.Phase = Pause
	}
	return nil
}
