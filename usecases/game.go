package usecases

import (
	"coinche/domain"
)

type GameUsecasesInterface interface {
	ListGames() ([]domain.Game, error)
	GetGame(gameID int) (domain.Game, error)
	CreateGame(name string) int
	JoinGame(gameID int, playerName string) (domain.Game, error)
	LeaveGame(gameID int, playerName string) error
}

type GameRepositoryInterface interface {
	ListGames() ([]domain.Game, error)
	GetGame(gameID int) (domain.Game, error)
	CreateGame(game domain.Game) (int, error)
	UpdatePlayer(gameID int, playerName string, players domain.Player) error
	UpdateGame(game domain.Game) error
}

type GameUsecases struct {
	GameUsecasesInterface
	Repo GameRepositoryInterface
}

func (s *GameUsecases) ListGames() ([]domain.Game, error) {
	return s.Repo.ListGames()
}

func (s *GameUsecases) GetGame(gameID int) (domain.Game, error) {
	return s.Repo.GetGame(gameID)
}

func (s *GameUsecases) CreateGame(name string) (int, error) {
	game := domain.NewGame(name)
	return s.Repo.CreateGame(game)
}

func (s *GameUsecases) JoinGame(gameID int, playerName string) (domain.Game, error) {
	game, err := s.Repo.GetGame(gameID)
	if err != nil {
		return domain.Game{}, err
	}
	err = game.AddPlayer(playerName)
	if err != nil {
		return domain.Game{}, err
	}
	err = s.Repo.UpdateGame(game)
	if err != nil {
		return domain.Game{}, err
	}

	game, err = s.Repo.GetGame(gameID)

	return game, err
}

func (s *GameUsecases) LeaveGame(gameID int, playerName string) error {
	game, err := s.Repo.GetGame(gameID)
	if err != nil {
		return err
	}

	err = game.RemovePlayer(playerName)
	if err != nil {
		return err
	}
	err = s.Repo.UpdateGame(game)
	if err != nil {
		return err
	}
	game, err = s.Repo.GetGame(gameID)

	return err
}

func (s *GameUsecases) JoinTeam(gameID int, playerName string, teamName string) error {
	game, err := s.Repo.GetGame(gameID)
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

func (s *GameUsecases) LeaveTeam(gameID int, playerName string) error {
	game, err := s.Repo.GetGame(gameID)
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

func (s *GameUsecases) StartGame(gameID int) error {
	game, err := s.Repo.GetGame(gameID)
	if err != nil {
		return err
	}
	err = game.StartBidding()
	if err != nil {
		return err
	}
	err = s.Repo.UpdateGame(game)
	return err
}

func (s *GameUsecases) Bid(gameID int, playerName string, value domain.BidValue, color domain.Color) error {
	game, err := s.Repo.GetGame(gameID)
	if err != nil {
		return err
	}

	err = game.PlaceBid(playerName, value, color)
	if err != nil {
		return err
	}
	err = s.Repo.UpdateGame(game)
	return err
}

func (s *GameUsecases) Pass(gameID int, playerName string) error {
	game, err := s.Repo.GetGame(gameID)
	if err != nil {
		return err
	}

	err = game.Pass(playerName)
	if err != nil {
		return err
	}

	err = s.Repo.UpdateGame(game)
	return err
}

func (s *GameUsecases) Coinche(gameID int, playerName string) error {
	game, err := s.Repo.GetGame(gameID)
	if err != nil {
		return err
	}

	err = game.Coinche(playerName)
	if err != nil {
		return err
	}
	err = s.Repo.UpdateGame(game)
	return err
}

func (s *GameUsecases) PlayCard(gameID int, playerName string, card domain.CardID) error {
	game, err := s.Repo.GetGame(gameID)
	if err != nil {
		return err
	}

	err = game.Play(playerName, card)
	if err != nil {
		return err
	}

	err = s.Repo.UpdateGame(game)
	return err
}

func NewGameUsecases(repository GameRepositoryInterface) *GameUsecases {
	return &GameUsecases{Repo: repository}
}
