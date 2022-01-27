package adapters

import (
	"coinche/app"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

func (s *dbGameService) ListGames() []app.Game {
	var games []app.Game
	err := s.db.Select(&games, "SELECT * FROM game ")
	if err != nil {
		fmt.Println(err)
	}

	return games
}
