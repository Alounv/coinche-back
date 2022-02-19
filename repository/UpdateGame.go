package repository

import (
	"coinche/domain"
)

func (s *GameRepository) UpdateGame(gameID int, phase domain.Phase) error {
	_, err := s.db.Exec(
		`
		UPDATE game
		SET phase = $1
		WHERE id = $2
		`,
		phase,
		gameID,
	)
	if err != nil {
		return err
	}
	return nil
}
