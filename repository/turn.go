package repository

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var turnSchema = `
CREATE TABLE turn (
	id serial PRIMARY KEY NOT NULL,
  position integer NOT NULL,
	gameid integer NOT NULL REFERENCES game(id),
	winner  text NOT NULL,
	plays json NOT NULL DEFAULT '[]'
)`

func createTurn(turn domain.Turn, tx *sqlx.Tx, gameID int, position int) error {
	var plays []byte
	plays, err := json.Marshal(turn.Plays)
	utilities.PanicIfErr(err)
	_, err = tx.Exec(
		`
			INSERT INTO turn (gameid, winner, plays, position)
			VALUES ($1, $2, $3, $4)
			`,
		gameID,
		turn.Winner,
		plays,
		position,
	)
	return err
}

func createTurns(tx *sqlx.Tx, gameID int, turns []domain.Turn) error {
	for position, turn := range turns {
		err := createTurn(turn, tx, gameID, position)
		if err != nil {
			return err
		}
	}
	return nil
}

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

func getTurns(tx *sqlx.Tx, gameID int) ([]domain.Turn, error) {
	turns := []domain.Turn{}

	type DBTurn struct {
		Winner string
		Plays  []byte
	}

	var dbTurns []DBTurn

	err := tx.Select(&dbTurns, `SELECT winner, plays FROM turn WHERE gameid=$1 ORDER BY position`, gameID)
	if err != nil {
		return turns, err
	}

	for _, dbTurn := range dbTurns {
		var plays []domain.Play
		err := json.Unmarshal(dbTurn.Plays, &plays)
		if err != nil {
			return turns, err
		}

		turns = append(turns, domain.Turn{
			Winner: dbTurn.Winner,
			Plays:  plays,
		})

	}

	return turns, nil
}
