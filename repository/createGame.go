package repository

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"
)

func (s *GameRepository) CreateGame(game domain.Game) (int, error) {
	tx := s.db.MustBegin()

	var gameID int

	deck, err := json.Marshal(game.Deck)
	utilities.PanicIfErr(err)

	err = tx.QueryRow(
		`
		INSERT INTO game (name, phase, deck) 
		VALUES ($1, $2, $3) 
		RETURNING id
		`,
		game.Name,
		game.Phase,
		deck,
	).Scan(&gameID)
	if err != nil {
		return 0, err
	}

	for playerName, player := range game.Players {
		_, err := tx.Exec(
			`
			INSERT INTO player (name, team, gameid) 
			VALUES ($1, $2, $3)
			`,
			playerName,
			player.Team,
			gameID,
		)
		if err != nil {
			return 0, err
		}
	}

	return gameID, tx.Commit()
}
