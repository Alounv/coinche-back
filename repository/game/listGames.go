package gamerepo

import (
	"coinche/domain"
	"fmt"

	"github.com/lib/pq"
)

func (s *GameRepositary) ListGames() []domain.Game {
	var games []domain.Game
	rows, err := s.db.Query("SELECT * FROM game ")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var game domain.Game
		err = rows.Scan(
			&game.ID,
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
