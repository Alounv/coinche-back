package gamerepo

import (
	"coinche/domain"
	"strings"
)

func (s *GameRepository) UpdatePlayers(id int, players []string, phase domain.Phase) error {
	_, err := s.db.Exec(
		`
		UPDATE game
		SET players = string_to_array($1, ','), phase = $2
		WHERE id = $3
		`,
		strings.Join(players, ","),
		phase,
		id,
	)
	return err
}
