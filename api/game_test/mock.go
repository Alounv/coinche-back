package gameapitest

import (
	"coinche/domain"
	"errors"
	"sort"
)

type MockGameUsecases struct {
	games    map[int]domain.Game
	setCalls []string
}

func (s *MockGameUsecases) GetGame(id int) (domain.Game, error) {
	game, existingID := s.games[id]
	if !existingID {
		return domain.Game{}, errors.New("NO GAME FOUND")
	}
	return game, nil
}

func (s *MockGameUsecases) CreateGame(name string, creatorName string) int {
	s.setCalls = append(s.setCalls, name)
	return 1
}

type ByID []domain.Game

func (a ByID) Len() int           { return len(a) }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (s *MockGameUsecases) ListGames() []domain.Game {
	var games []domain.Game
	for id, val := range s.games {
		val.ID = id
		games = append(games, val)
	}
	sort.Sort(ByID(games))
	return games
}

func (s *MockGameUsecases) JoinGame(id int, playerName string) (domain.Game, error) {
	var err error
	game, existingID := s.games[id]
	if !existingID {
		err = errors.New("TEST JOIN FAIL")
	}

	gameWithNewPlayer := domain.Game{
		ID:        id,
		Name:      game.Name,
		CreatedAt: game.CreatedAt,
		Players:   append(game.Players, playerName),
	}
	return gameWithNewPlayer, err
}

func (s *MockGameUsecases) LeaveGame(id int, playerName string) error {
	var err error
	_, existingID := s.games[id]
	if !existingID {
		err = errors.New("TEST LEAVE FAIL")
	}

	return err
}
