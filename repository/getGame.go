package repository

import (
	"coinche/domain"
)

func (s *GameRepository) GetGame(gameID int) (domain.Game, error) {
	tx := s.db.MustBegin()

	var game domain.Game
	err := tx.QueryRow(`SELECT * FROM game WHERE id=$1`, gameID).Scan(
		&game.ID,
		&game.Name,
		&game.CreatedAt,
		&game.Phase,
	)
	if err != nil {
		return domain.Game{}, err
	}

	rows, err := tx.Query(`SELECT name, team FROM player WHERE gameid=$1`, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Players = map[string]domain.Player{}

	for rows.Next() {
		var playerName string
		var teamName string
		err := rows.Scan(&playerName, &teamName)
		if err != nil {
			return domain.Game{}, err
		}

		game.Players[playerName] = domain.Player{Team: teamName}
	}

	return game, tx.Commit()
}
