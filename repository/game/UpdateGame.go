package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) UpdateGame(id int, phase domain.Phase) error {
	_, err := s.db.Exec(
		`
		UPDATE game
		SET phase = $1
		WHERE id = $2
		`,
		phase,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
