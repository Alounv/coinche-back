package usecases

import (
	"coinche/domain"
	"fmt"
)

type GameUsecaseInterface interface {
	ListGames() []domain.Game
	GetGame(id int) domain.Game
	CreateGame(name string) int
	JoinGame(id int, playerName string) error
}

type GameRepositoryInterface interface {
	ListGames() []domain.Game
	GetGame(id int) domain.Game
	CreateGame(name string) int
	UpdateGame(id int, playerNames []string) error
}

type GameService struct {
	GameUsecaseInterface
	Repo GameRepositoryInterface
}

func (s *GameService) ListGames() []domain.Game {
	return s.Repo.ListGames()
}

func (s *GameService) GetGame(id int) domain.Game {
	return s.Repo.GetGame(id)
}

func (s *GameService) CreateGame(name string) int {
	return s.Repo.CreateGame(name)
}

func (s *GameService) JoinGame(id int, playerName string) error {
	playersNames := s.Repo.GetGame(id).Players
	fmt.Print(playersNames)
	playersNames = append(playersNames, playerName)
	fmt.Print(playersNames)
	return s.Repo.UpdateGame(id, playersNames)
}

func NewGameService(repository GameRepositoryInterface) *GameService {
	return &GameService{Repo: repository}
}
