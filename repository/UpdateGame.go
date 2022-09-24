package repository

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func getCurrentTurnsCount(tx *sqlx.Tx, gameID int) (int, error) {
	currentTurnsCount := 0
	row := tx.QueryRow(`SELECT COUNT (*) FROM turn WHERE gameid=$1`, gameID)
	err := row.Scan(&currentTurnsCount)
	fmt.Println("GET COUNT", currentTurnsCount, err)
	return currentTurnsCount, err
}

func updateTurn(turn domain.Turn, tx *sqlx.Tx, gameID int, position int) error {
	var plays []byte
	plays, err := json.Marshal(turn.Plays)
	utilities.PanicIfErr(err)
	_, err = tx.Exec(
		`
			UPDATE turn
      SET winner = $2, plays = $3 
			WHERE position = $4 AND gameid = $1
			`,
		gameID,
		turn.Winner,
		plays,
		position,
	)
	return err
}

func updateTurns(tx *sqlx.Tx, gameID int, turns []domain.Turn) error {
	count, err := getCurrentTurnsCount(tx, gameID)
	if err != nil {
		return err
	}

	for index, turn := range turns[count:] {
		err := createTurn(turn, tx, gameID, index+count)
		if err != nil {
			return err
		}
	}

	for position, turn := range turns[:count] {
		err := updateTurn(turn, tx, gameID, position)
		if err != nil {
			return err
		}
	}

	return nil
}

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
	return tx.Commit()
}
