package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) ListGames() ([]domain.Game, error) {
	var games []domain.Game

	rows, err := s.db.Query("SELECT * FROM game ")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var game domain.Game
		err = rows.Scan(
			&game.ID,
			&game.Name,
			&game.CreatedAt,
			&game.Phase,
		)
		if err != nil {
			return nil, err
		}

		rows, err := s.db.Query("SELECT name FROM player WHERE gameid=$1", game.ID)
		if err != nil {
			return nil, err
		}

		game.Players = []string{}
		for rows.Next() {
			var playerName string
			err = rows.Scan(&playerName)
			game.Players = append(game.Players, playerName)
			if err != nil {
				return nil, err
			}
		}

		games = append(games, game)
	}

	return games, nil
}
