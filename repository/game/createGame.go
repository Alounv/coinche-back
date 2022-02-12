package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) CreateGame(game domain.Game) (int, error) {
	tx := s.db.MustBegin()

	var id int
	err := tx.QueryRow(
		`
		INSERT INTO game (name) 
		VALUES ($1) 
		RETURNING id
		`,
		game.Name,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, playerName := range game.Players {
		_, err := tx.Exec(
			`
			INSERT INTO player (name, gameid) 
			VALUES ($1, $2)
			`,
			playerName,
			id,
		)
		if err != nil {
			return 0, err
		}
	}

	return id, tx.Commit()
}
