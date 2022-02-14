package usecases

import (
	"coinche/domain"
)

type GameUsecasesInterface interface {
	ListGames() ([]domain.Game, error)
	GetGame(id int) (domain.Game, error)
	CreateGame(name string) int
	JoinGame(id int, playerName string) (domain.Game, error)
	LeaveGame(id int, playerName string) error
}

type GameRepositoryInterface interface {
	ListGames() ([]domain.Game, error)
	GetGame(id int) (domain.Game, error)
	CreateGame(game domain.Game) (int, error)
	UpdatePlayers(id int, players map[string]domain.Player, phase domain.Phase) error
	UpdatePlayer(id int, playerName string, players domain.Player) error
}

type GameUsecases struct {
	GameUsecasesInterface
	Repo GameRepositoryInterface
}

func (s *GameUsecases) ListGames() ([]domain.Game, error) {
	return s.Repo.ListGames()
}

func (s *GameUsecases) GetGame(id int) (domain.Game, error) {
	return s.Repo.GetGame(id)
}

func (s *GameUsecases) CreateGame(name string) (int, error) {
	game := domain.NewGame(name)
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

func (s *GameUsecases) JoinTeam(id int, playerName string, teamName string) error {
	game, err := s.Repo.GetGame(id)
	if err != nil {
		return err
	}
	err = game.AssignTeam(playerName, teamName)
	if err != nil {
		return err
	}
	err = s.Repo.UpdatePlayer(game.ID, playerName, game.Players[playerName])
	return err
}

func (s *GameUsecases) LeaveTeam(id int, playerName string) error {
	game, err := s.Repo.GetGame(id)
	if err != nil {
		return err
	}
	err = game.ClearTeam(playerName)
	if err != nil {
		return err
	}
	err = s.Repo.UpdatePlayer(game.ID, playerName, game.Players[playerName])
	return err
}

func NewGameUsecases(repository GameRepositoryInterface) *GameUsecases {
	return &GameUsecases{Repo: repository}
}
