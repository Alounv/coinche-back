package adapters

import (
	"coinche/app"

	_ "github.com/jackc/pgx/stdlib"
)

type GameService struct {
	DbGameService *dbGameService
}

func (s *GameService) ListGames() []app.Game {
	return s.DbGameService.ListGames()
}

func (s *GameService) GetGame(id int) app.Game {
	return s.DbGameService.GetGame(id)
}

func (s *GameService) CreateGame(name string) int {
	return s.DbGameService.CreateGame(name)
}

func (s *GameService) JoinGame(id int, playerName string) error {
	return nil
}

func NewGameService(dsn string) *GameService {
	dbService := NewDBGameService(dsn)
	return &GameService{dbService}
}