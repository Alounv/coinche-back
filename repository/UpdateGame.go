package repository

import (
	"coinche/domain"
	"encoding/json"
)

func (r *GameRepository) UpdateGame(game domain.Game) error {
	tx := r.db.MustBegin()

	deck, err := json.Marshal(game.Deck)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`
		UPDATE game
		SET phase = $2, Deck = $3 
		WHERE id = $1
		`,
		game.ID,
		game.Phase,
		deck,
	)

	if err != nil {
		return err
	}

	err = createBids(tx, game.ID, game.Bids)
	if err != nil {
		return err
	}

	err = updateTurns(tx, game.ID, game.Turns)
	if err != nil {
		return err
	}

	err = updatePlayers(tx, game.ID, game.Players)
	if err != nil {
		return err
	}

	err = createPoints(tx, game.ID, game.Points)
	if err != nil {
		return err
	}

	err = createScores(tx, game.ID, game.Scores)
	if err != nil {
		return err
	}

	return tx.Commit()
}
