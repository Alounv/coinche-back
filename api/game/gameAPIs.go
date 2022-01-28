package gameapi

import (
	"coinche/domain"
)

type GameUsecase interface {
	ListGames() []domain.Game
	GetGame(id int) domain.Game
	CreateGame(name string) int
	JoinGame(id int, playerName string) error
}

type GameAPIs struct {
	GameService GameUsecase
}
