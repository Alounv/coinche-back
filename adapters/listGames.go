package adapters

import (
	"coinche/app"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
)

func (s *dbGameService) ListGames() []app.Game {
	var games []app.Game
	rows, err := s.db.Query("SELECT * FROM game ")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var game app.Game
		err = rows.Scan(
			&game.Id,
			&game.Name,
			&game.CreatedAt,
			(*pq.StringArray)(&game.Players),
		)
		if err != nil {
			fmt.Println(err)
		}
		games = append(games, game)
	}

	return games
}
