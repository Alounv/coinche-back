package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var pointSchema = `
CREATE TABLE IF NOT EXISTS point (
	id serial PRIMARY KEY NOT NULL,
	gameid integer NOT NULL REFERENCES game(id),
	team  text NOT NULL,
	value integer NOT NULL
)`

var scoreSchema = `
CREATE TABLE IF NOT EXISTS score (
	id serial PRIMARY KEY NOT NULL,
	gameid integer NOT NULL REFERENCES game(id),
	team  text NOT NULL,
	value integer NOT NULL
)`

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

func updatePointsOrScore(tx *sqlx.Tx, gameID int, collection string, team string, value int) error {
	query := fmt.Sprintf(`
      UPDATE %s
      SET value = $3
      WHERE gameid = $1 AND team = $2
			`, collection)

	_, err := tx.Exec(
		query,
		gameID,
		team,
		value,
	)
	return err
}

func createAndUpdatePointsOrScores(tx *sqlx.Tx, gameID int, collection string, current map[string]int, values map[string]int) error {
	for team, value := range values {
		shouldCreate := false
		if _, ok := current[team]; !ok {
			shouldCreate = true
		}

		if shouldCreate {
			err := createPointsOrScore(tx, gameID, collection, team, value)
			if err != nil {
				return err
			}
		} else {
			err := updatePointsOrScore(tx, gameID, collection, team, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createAndUpdateScores(tx *sqlx.Tx, gameID int, scores map[string]int) error {
	currentScores, err := getScoresOrPoints(tx, gameID, "score")
	if err != nil {
		return err
	}

	err = createAndUpdatePointsOrScores(tx, gameID, "score", currentScores, scores)
	if err != nil {
		return err
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

func getScoresOrPoints(tx *sqlx.Tx, gameID int, collection string) (map[string]int, error) {
	type TeamValue struct {
		Team  string
		Value int
	}

	data := []TeamValue{}

	query := fmt.Sprintf(`SELECT team, value FROM %s WHERE gameid = $1`, collection)
	err := tx.Select(&data, query, gameID)

	if err != nil {
		return nil, err
	}

	var scores map[string]int = map[string]int{}
	for _, d := range data {
		scores[d.Team] = d.Value
	}

	return scores, nil
}
