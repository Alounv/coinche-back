package usecases

import (
	"coinche/domain"
)

type GameUsecasesInterface interface {
	ListGames() []domain.Game
	GetGame(id int) (domain.Game, error)
	CreateGame(name string, creatorName string) int
	JoinGame(id int, playerName string) (domain.Game, error)
	LeaveGame(id int, playerName string) error
}

type GameRepositoryInterface interface {
	ListGames() []domain.Game
	GetGame(id int) (domain.Game, error)
	CreateGame(game domain.Game) int
	UpdatePlayers(id int, players []string, phase domain.Phase) error
}

type GameUsecases struct {
	GameUsecasesInterface
	Repo GameRepositoryInterface
}

func (s *GameUsecases) ListGames() []domain.Game {
	return s.Repo.ListGames()
}

func (s *GameUsecases) GetGame(id int) (domain.Game, error) {
	return s.Repo.GetGame(id)
}

func (s *GameUsecases) CreateGame(name string, creatorName string) int {
	game := domain.NewGame(name, creatorName)
	return s.Repo.CreateGame(game)
}

func (s *GameUsecases) JoinGame(id int, playerName string) (domain.Game, error) {
	game, err := s.Repo.GetGame(id)
	if err != nil {
		return game, err
	}
	err = game.AddPlayer(playerName)
	if err != nil {
		return game, err
	}
	err = s.Repo.UpdatePlayers(game.ID, game.Players, game.Phase)
	return game, err
}

func (s *GameUsecases) LeaveGame(id int, playerName string) error {
	game, err := s.Repo.GetGame(id)
	if err != nil {
		return err
	}
	err = game.RemovePlayer(playerName)
	if err != nil {
		return err
	}
	err = s.Repo.UpdatePlayers(game.ID, game.Players, game.Phase)
	return err
}

func NewGameUsecases(repository GameRepositoryInterface) *GameUsecases {
	return &GameUsecases{Repo: repository}
}
