package repository

import (
	"coinche/domain"
)

func (s *GameRepository) CreateGames(games []domain.Game) error {
	tx := s.db.MustBegin()

	for _, game := range games {

		_, err := s.createAGame(game, tx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
