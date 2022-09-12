package usecases

import (
	"coinche/domain"
	"errors"
)

type MockGameRepo struct {
	games         map[int]domain.Game
	creationCalls int
}

func (repo *MockGameRepo) ListGames() ([]domain.Game, error) {
	var games []domain.Game
	for gameID, val := range repo.games {
		val.ID = gameID
		games = append(games, val)
	}
	return games, nil
}

func (repo *MockGameRepo) GetGame(gameID int) (domain.Game, error) {
	game, ok := repo.games[gameID]
	if !ok {
		return domain.Game{}, errors.New("GAME NOT FOUND")
	}
	return game, nil
}

func (repo *MockGameRepo) CreateGame(game domain.Game) (int, error) {
	gameID := len(repo.games)
	repo.creationCalls = repo.creationCalls + 1
	return gameID, nil
}

func (repo *MockGameRepo) UpdatePlayers(gameID int, players map[string]domain.Player, phase domain.Phase) error {
	game, ok := repo.games[gameID]
	if !ok {
		return errors.New("GAME NOT FOUND")
	}
	game.Players = players
	game.Phase = phase
	repo.games[gameID] = game
	return nil
}

func (repo *MockGameRepo) UpdatePlayer(gameID int, playerName string, player domain.Player) error {
	game, ok := repo.games[gameID]
	if !ok {
		return errors.New("GAME NOT FOUND")
	}
	game.Players[playerName] = player
	repo.games[gameID] = game
	return nil
}
func (repo *MockGameRepo) UpdateGame(gameID int, phase domain.Phase) error {
	game, ok := repo.games[gameID]
	if !ok {
		return errors.New("GAME NOT FOUND")

	}
	game.Phase = phase
	repo.games[gameID] = game
	return nil
}

func NewMockGameRepo(games map[int]domain.Game) MockGameRepo {
	return MockGameRepo{
		games:         games,
		creationCalls: 0,
	}
}
