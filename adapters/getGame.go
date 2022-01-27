package adapters

import (
	"coinche/app"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

func (s *dbGameService) GetGame(id int) app.Game {
	var game app.Game
	err := s.db.Get(&game, "SELECT * FROM game WHERE id=$1", id)
	if err != nil {
		fmt.Println(err)
	}

	return game
}
