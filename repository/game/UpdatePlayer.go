package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) UpdatePlayer(id int, playerName string, player domain.Player) error {
	_, err := s.db.Exec(
		`
		INSERT INTO player (name, team, gameid) 
		VALUES ($1, $2, $3)
		`,
		playerName,
		player.Team,
		id,
	)
	if err != nil {
		return err
	}

	return err
}
