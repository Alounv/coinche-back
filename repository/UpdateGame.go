package repository

import (
	"coinche/domain"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
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

	if game.Bids != nil && len(game.Bids) == 0 {
		err = resetItems(tx, game.ID, "bid")
		if err != nil {
			return err
		}
	} else {
		err = createBids(tx, game.ID, game.Bids)
		if err != nil {
			return err
		}
	}

	if game.Turns != nil && len(game.Turns) == 0 {
		err = resetItems(tx, game.ID, "turn")
		if err != nil {
			return err
		}
	} else {
		err = updateTurns(tx, game.ID, game.Turns)
		if err != nil {
			return err
		}
	}

	err = updatePlayers(tx, game.ID, game.Players)
	if err != nil {
		return err
	}

	if game.Points != nil && len(game.Points) == 0 {
		err = resetItems(tx, game.ID, "point")
		if err != nil {
			return err
		}
	} else {
		err = createPoints(tx, game.ID, game.Points)
		if err != nil {
			return err
		}
	}

	err = createAndUpdateScores(tx, game.ID, game.Scores)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func resetItems(tx *sqlx.Tx, gameID int, collection string) error {
	query := fmt.Sprint("DELETE FROM ", collection, " WHERE gameid = $1")
	_, err := tx.Exec(query, gameID)
	if err != nil {
		return err
	}
	return nil
}
