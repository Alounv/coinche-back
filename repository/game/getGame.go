package gamerepo

import (
	"coinche/domain"
	"fmt"

	"github.com/lib/pq"
)

func (s *GameRepo) GetGame(id int) domain.Game {
	var game domain.Game
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
