package usecases

import (
	"coinche/domain"
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
	return *repo.games[id], nil
}

func (repo *MockGameRepo) CreateGame(game domain.Game) int {
	newId := len(repo.games)
	repo.creationCalls = repo.creationCalls + 1
	return newId
}

func (repo *MockGameRepo) UpdatePlayers(id int, players []string) error {
	repo.games[id].Players = players
	return nil
}
