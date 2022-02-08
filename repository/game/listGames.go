package gamerepo

import (
	"coinche/domain"

	"github.com/lib/pq"
)

func (s *GameRepository) ListGames() []domain.Game {
	var games []domain.Game
	rows, err := s.db.Query("SELECT * FROM game ")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var game domain.Game
		err = rows.Scan(
			&game.ID,
			&game.Name,
			&game.CreatedAt,
			&game.Phase,
			(*pq.StringArray)(&game.Players),
		)
		if err != nil {
			panic(err)
		}
		games = append(games, game)
	}

	return games
}
