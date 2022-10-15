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
	game.ID = gameID
	if !ok {
		return domain.Game{}, errors.New("GAME NOT FOUND")
	}

	return game, nil
}

func (repo *MockGameRepo) CreateGame(game domain.Game) (int, error) {
	var gameID int
	if game.ID == 0 {
		gameID = len(repo.games) + 1
	} else {
		gameID = game.ID
	}

	game.Root = gameID

	repo.creationCalls = repo.creationCalls + 1
	repo.games[gameID] = game
	return gameID, nil
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

func (repo *MockGameRepo) UpdateGame(game domain.Game) error {
	repoGame, ok := repo.games[game.ID]
	if !ok {
		return errors.New("GAME NOT FOUND")

	}
	repoGame.Phase = game.Phase
	repoGame.Bids = game.Bids
	repoGame.Deck = game.Deck
	repoGame.Players = game.Players
	repoGame.Turns = game.Turns
	repoGame.Points = game.Points
	repoGame.Scores = game.Scores
	repoGame.Root = game.Root

	repo.games[game.ID] = repoGame
	return nil
}

func (repo *MockGameRepo) DeleteGame(gameID int) error {
	delete(repo.games, gameID)
	return nil
}

func NewMockGameRepo(games map[int]domain.Game) MockGameRepo {
	for gameID, game := range games {
		game.Root = gameID
		games[gameID] = game
	}
	return MockGameRepo{
		games:         games,
		creationCalls: 0,
	}
}
