package repository

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func countDocumentsInGame(tx *sqlx.Tx, gameID int, collection string) (int, error) {
	count := 0
	query := fmt.Sprintf(`SELECT COUNT (*) FROM %s WHERE gameid=$1`, collection)
	row := tx.QueryRow(query, gameID)
	err := row.Scan(&count)
	return count, err
}

func getTurnsCount(tx *sqlx.Tx, gameID int) (int, error) {
	turnsCount, err := countDocumentsInGame(tx, gameID, "turn")
	return turnsCount, err
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

func createPointsOrScore(tx *sqlx.Tx, gameID int, collection string, team string, value int) error {
	query := fmt.Sprintf(`
			INSERT INTO %s (gameid, team, value)
			VALUES ($1, $2, $3)
			`, collection)

	_, err := tx.Exec(
		query,
		gameID,
		team,
		value,
	)
	return err
}

func updateTurns(tx *sqlx.Tx, gameID int, turns []domain.Turn) error {
	count, err := getTurnsCount(tx, gameID)
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

func createScores(tx *sqlx.Tx, gameID int, scores map[string]int) error {
	for team, teamScore := range scores {
		err := createPointsOrScore(tx, gameID, "score", team, teamScore)
		if err != nil {
			return err
		}
	}

	return nil
}

func createPoints(tx *sqlx.Tx, gameID int, points map[string]int) error {
	for team, teamPoints := range points {
		err := createPointsOrScore(tx, gameID, "point", team, teamPoints)
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
