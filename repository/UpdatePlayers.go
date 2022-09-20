package repository

import (
	"coinche/domain"

	"github.com/jmoiron/sqlx"
)

func getCurrentPlayers(tx *sqlx.Tx, gameID int) (map[string]int, error) {
	var currentPlayers map[string]int = map[string]int{}

	rows, err := tx.Query(`SELECT name FROM player WHERE gameid=$1`, gameID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var playerName string
		err := rows.Scan(&playerName)
		if err != nil {
			return nil, err
		}

		currentPlayers[playerName] = 0
	}

	return currentPlayers, nil
}

func deletePlayers(currentPlayers map[string]int, players map[string]domain.Player, gameID int, tx *sqlx.Tx) error {
	for currentPlayerName := range currentPlayers {
		shouldDelete := false
		if _, ok := players[currentPlayerName]; !ok {
			shouldDelete = true
		}

		if shouldDelete {
			_, err := tx.Exec(
				`
				DELETE FROM player
				WHERE name = $1 AND gameid = $2
				`,
				currentPlayerName,
				gameID,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createAndUpdatePlayers(currentPlayers map[string]int, players map[string]domain.Player, gameID int, tx *sqlx.Tx) error {
	for playerName, player := range players {
		shouldCreate := false
		if _, ok := currentPlayers[playerName]; !ok {
			shouldCreate = true
		}

		if shouldCreate {
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
				return err
			}
		} else {
			err := updatePlayer(tx, gameID, playerName, player)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updatePlayers(tx *sqlx.Tx, gameID int, players map[string]domain.Player) error {
	currentPlayers, err := getCurrentPlayers(tx, gameID)
	if err != nil {
		return err
	}

	err = deletePlayers(currentPlayers, players, gameID, tx)
	if err != nil {
		return err
	}

	err = createAndUpdatePlayers(currentPlayers, players, gameID, tx)
	if err != nil {
		return err
	}

	return nil
}
