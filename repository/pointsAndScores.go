package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var pointSchema = `
CREATE TABLE point (
	id serial PRIMARY KEY NOT NULL,
	gameid integer NOT NULL REFERENCES game(id),
	team  text NOT NULL,
	value integer NOT NULL
)`

var scoreSchema = `
CREATE TABLE score (
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
