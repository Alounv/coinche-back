package gamerepo

import (
	"coinche/domain"

	"github.com/lib/pq"
)

func (s *GameRepository) GetGame(id int) (domain.Game, error) {
	var game domain.Game
	err := s.db.QueryRow(`SELECT * FROM game WHERE id=$1`, id).Scan(
		&game.ID,
		&game.Name,
		&game.CreatedAt,
		(*pq.StringArray)(&game.Players),
	)

	return game, err
}
