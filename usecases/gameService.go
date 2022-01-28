package usecases

import (
	"coinche/domain"
	"fmt"
)

type GameUsecasesInterface interface {
	ListGames() []domain.Game
	GetGame(id int) domain.Game
	CreateGame(name string, creatorName string) int
	JoinGame(id int, playerName string) error
}

type GameRepositoryInterface interface {
	ListGames() []domain.Game
	GetGame(id int) domain.Game
	CreateGame(game domain.Game) int
	UpdateGame(id int, playerNames []string) error
}

type GameUsecases struct {
	GameUsecasesInterface
	Repo GameRepositoryInterface
}

func (s *GameUsecases) ListGames() []domain.Game {
	return s.Repo.ListGames()
}

func (s *GameUsecases) GetGame(id int) domain.Game {
	return s.Repo.GetGame(id)
}

func (s *GameUsecases) CreateGame(name string, creatorName string) int {
	game := domain.NewGame(name, creatorName)
	return s.Repo.CreateGame(game)
}

func (s *GameUsecases) JoinGame(id int, playerName string) error {
	playersNames := s.Repo.GetGame(id).Players
	fmt.Print(playersNames)
	playersNames = append(playersNames, playerName)
	fmt.Print(playersNames)
	return s.Repo.UpdateGame(id, playersNames)
}

func NewGameUsecases(repository GameRepositoryInterface) *GameUsecases {
	return &GameUsecases{Repo: repository}
}
