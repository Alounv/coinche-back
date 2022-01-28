package app

import (
	"coinche/domain"
	gamerepo "coinche/repository/game"
	"fmt"
)

type GameService struct {
	GameRepo *gamerepo.GameRepo
}

func (s *GameService) ListGames() []domain.Game {
	return s.GameRepo.ListGames()
}

func (s *GameService) GetGame(id int) domain.Game {
	return s.GameRepo.GetGame(id)
}

func (s *GameService) CreateGame(name string) int {
	return s.GameRepo.CreateGame(name)
}

func (s *GameService) JoinGame(id int, playerName string) error {
	playersNames := s.GameRepo.GetGame(id).Players
	fmt.Print(playersNames)
	playersNames = append(playersNames, playerName)
	fmt.Print(playersNames)
	return s.GameRepo.UpdateGame(id, playersNames)
}

func NewGameService(dsn string) *GameService {
	dbService := gamerepo.NewGameRepo(dsn)
	return &GameService{dbService}
}