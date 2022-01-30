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
	return s.games[id], nil
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

func (s *MockGameUsecases) JoinGame(id int, playerName string) error {
	var err error
	if _, existingID := s.games[id]; !existingID {
		err = errors.New("GAME NOT FOUND")
	} else if s.games[id].IsFull() {
		err = errors.New("GAME IS FULL")
	}

	return err
}
