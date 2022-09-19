package repository

import (
	"coinche/domain"
)

func (s *GameRepository) UpdateGame(game domain.Game) error {
	tx := s.db.MustBegin()

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

	err = s.CreateBids(tx, game.ID, game.Bids)

	if err != nil {
		return err
	}
	return tx.Commit()
}
