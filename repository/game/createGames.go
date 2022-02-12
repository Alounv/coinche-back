package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) CreateGames(games []domain.Game) error {
	tx := s.db.MustBegin()

	for _, game := range games {
		_, err := tx.Exec(
			`
			INSERT INTO game (id, name, createdAt, phase)
			VALUES ($1, $2, $3, $4)
			`,
			game.ID,
			game.Name,
			game.CreatedAt,
			game.Phase,
		)
		if err != nil {
			return err
		}

		for _, playerName := range game.Players {
			_, err := tx.Exec(
				`
				INSERT INTO player (name, gameid) 
				VALUES ($1, $2)
				`,
				playerName,
				game.ID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
