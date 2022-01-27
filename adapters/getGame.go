package adapters

import (
	"coinche/app"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
)

func (s *dbGameService) GetGame(id int) app.Game {
	var game app.Game
	err := s.db.QueryRow(`SELECT * FROM game WHERE id=$1`, id).Scan(
		&game.Id,
		&game.Name,
		&game.CreatedAt,
		(*pq.StringArray)(&game.Players),
	)
	if err != nil {
		fmt.Println(err)
	}

	return game
}
