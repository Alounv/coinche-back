package gameapi

import (
	"coinche/domain"
	"errors"
	"sort"
)

type MockGameService struct {
	games    map[int]domain.Game
	setCalls []string
}

func (s *MockGameService) GetGame(id int) domain.Game {
	return s.games[id]
}

func (s *MockGameService) CreateGame(name string) int {
	s.setCalls = append(s.setCalls, name)
	return 1
}

type ById []domain.Game

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (s *MockGameService) ListGames() []domain.Game {
	var games []domain.Game
	for _, val := range s.games {
		games = append(games, val)
	}
	sort.Sort(ById(games))
	return games
}

func (s *MockGameService) JoinGame(id int, playerName string) error {
	var err error
	if _, existingId := s.games[id]; !existingId {
		err = errors.New("Game not found")
	} else if s.games[id].Full {
		err = errors.New("Game is full")
	}

	return err
}
