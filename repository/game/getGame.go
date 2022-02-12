package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) GetGame(id int) (domain.Game, error) {
	tx := s.db.MustBegin()

	var game domain.Game
	err := tx.QueryRow(`SELECT * FROM game WHERE id=$1`, id).Scan(
		&game.ID,
		&game.Name,
		&game.CreatedAt,
		&game.Phase,
	)
	if err != nil {
		return domain.Game{}, err
	}

	rows, err := tx.Query(`SELECT name FROM player WHERE gameid=$1`, id)
	if err != nil {
		return domain.Game{}, err
	}

	game.Players = []string{}

	for rows.Next() {
		var playerName string
		err := rows.Scan(&playerName)
		if err != nil {
			return domain.Game{}, err
		}

		game.Players = append(game.Players, playerName)
	}

	return game, tx.Commit()
}
