package gameApi

import (
	"coinche/domain"
	"errors"
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

func (s *MockGameService) ListGames() []domain.Game {
	var games []domain.Game
	for _, val := range s.games {
		games = append(games, val)
	}
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
