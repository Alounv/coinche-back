package repository

import (
	"coinche/domain"
)

func (s *GameRepository) UpdateGame(game domain.Game) error {
	_, err := s.db.Exec(
		`
		UPDATE game
		SET phase = $1
		WHERE id = $2
		`,
		game.Phase,
		game.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
