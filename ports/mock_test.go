package ports

import (
	"coinche/app"
)

type MockGameService struct {
	names    map[int]string
	setCalls []string
}

func (s *MockGameService) GetGame(id int) app.Game {
	game := app.Game{Name: s.names[id], Id: id}
	return game
}

func (s *MockGameService) CreateGame(name string) int {
	s.setCalls = append(s.setCalls, name)
	return 1
}

func (s *MockGameService) ListGames() []app.Game {
	var games []app.Game
	for id, name := range s.names {
		games = append(games, app.Game{Name: name, Id: id})
	}
	return games
}
