package usecases

import (
	"coinche/domain"
	"errors"
)

type MockGameRepo struct {
	games         map[int]*domain.Game
	creationCalls int
}

func (repo *MockGameRepo) ListGames() []domain.Game {
	var games []domain.Game
	for id, val := range repo.games {
		val.ID = id
		games = append(games, *val)
	}
	return games
}

func (repo *MockGameRepo) GetGame(id int) (domain.Game, error) {
	game, ok := repo.games[id]
	if !ok {
		return domain.Game{}, errors.New("GAME NOT FOUND")
	}
	return *game, nil
}

func (repo *MockGameRepo) CreateGame(game domain.Game) int {
	newId := len(repo.games)
	repo.creationCalls = repo.creationCalls + 1
	return newId
}

func (repo *MockGameRepo) UpdatePlayers(id int, players []string, phase domain.Phase) error {
	game, ok := repo.games[id]
	if !ok {
		return errors.New("GAME NOT FOUND")
	}
	game.Players = players
	if game.IsFull() && (game.Phase == domain.Preparation) {
		game.Phase = domain.Teaming
	}
	return nil
}

func NewMockGameRepo(games map[int]*domain.Game) MockGameRepo {
	return MockGameRepo{
		games:         games,
		creationCalls: 0,
	}
}
