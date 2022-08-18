package repository

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"
)

func (s *GameRepository) CreateGames(games []domain.Game) error {
	tx := s.db.MustBegin()

	for _, game := range games {

		deck, err := json.Marshal(game.Deck)
		utilities.PanicIfErr(err)

		_, err = tx.Exec(
			`
			INSERT INTO game (id, name, createdAt, phase, deck)
			VALUES ($1, $2, $3, $4, $5)
			`,
			game.ID,
			game.Name,
			game.CreatedAt,
			game.Phase,
			deck,
		)
		if err != nil {
			return err
		}

		for playerName, player := range game.Players {
			_, err := tx.Exec(
				`
				INSERT INTO player (name, team, gameid) 
				VALUES ($1, $2, $3)
				`,
				playerName,
				player.Team,
				game.ID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
