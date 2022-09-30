package repository

import (
	"coinche/domain"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

var playerSchema = `
CREATE TABLE IF NOT EXISTS player (
	id serial PRIMARY KEY NOT NULL,
	name text NOT NULL,
	team text,
	gameid integer NOT NULL REFERENCES game(id),
	createdAt timestamp NOT NULL DEFAULT now(),
	initialOrder integer DEFAULT 0,
	cOrder integer DEFAULT 0,
	hand json NOT NULL DEFAULT '[]'
)`

func updatePlayer(tx *sqlx.Tx, gameID int, playerName string, player domain.Player) error {
	hand, err := json.Marshal(player.Hand)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
    UPDATE player
    SET gameid =$1, name = $2, team = $3, initialOrder = $4, cOrder = $5, hand = $6
    WHERE gameid = $1 AND name = $2
    `,
		gameID,
		playerName,
		player.Team,
		player.InitialOrder,
		player.Order,
		hand,
	)
	if err != nil {
		return err
	}

	return err
}

func (s *GameRepository) UpdatePlayer(gameID int, playerName string, player domain.Player) error {
	tx := s.db.MustBegin()
	err := updatePlayer(tx, gameID, playerName, player)
	if err != nil {
		return err
	}

	return tx.Commit()
}

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

func deletePlayer(tx *sqlx.Tx, gameID int, playerName string) error {
	_, err := tx.Exec(`DELETE FROM player WHERE gameid=$1 AND name=$2`, gameID, playerName)
	return err
}

func deletePlayers(currentPlayers map[string]int, players map[string]domain.Player, gameID int, tx *sqlx.Tx) error {
	for name := range currentPlayers {
		shouldDelete := false
		if _, ok := players[name]; !ok {
			shouldDelete = true
		}

		if shouldDelete {
			err := deletePlayer(tx, gameID, name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createPlayer(tx *sqlx.Tx, gameID int, playerName string, team string) error {
	_, err := tx.Exec(`INSERT INTO player (gameid, name, team) VALUES ($1, $2, $3)`,
		gameID,
		playerName,
		team,
	)
	return err
}

func createAndUpdatePlayers(currentPlayers map[string]int, players map[string]domain.Player, gameID int, tx *sqlx.Tx) error {
	for playerName, player := range players {
		shouldCreate := false
		if _, ok := currentPlayers[playerName]; !ok {
			shouldCreate = true
		}

		if shouldCreate {
			err := createPlayer(tx, gameID, playerName, player.Team)
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

func createFullPlayerForTest(tx *sqlx.Tx, gameID int, name string, player domain.Player) error {
	hand, err := json.Marshal(player.Hand)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`
			INSERT INTO player (name, team, gameid, initialOrder, cOrder, hand) 
			VALUES ($1, $2, $3, $4, $5, $6)
			`,
		name,
		player.Team,
		gameID,
		player.InitialOrder,
		player.Order,
		hand,
	)

	return err
}

func getPlayers(tx *sqlx.Tx, gameID int) (map[string]domain.Player, error) {
	type DBplayer struct {
		Name         string
		Team         string
		InitialOrder int
		COrder       int
		Hand         []byte
	}

	var dbPlayers []DBplayer
	var players map[string]domain.Player = map[string]domain.Player{}

	err := tx.Select(&dbPlayers, `SELECT name, team, initialOrder, cOrder, hand FROM player WHERE gameid=$1`, gameID)
	if err != nil {
		return players, err
	}

	for _, dbPlayer := range dbPlayers {
		var hand []domain.CardID
		err = json.Unmarshal(dbPlayer.Hand, &hand)
		if err != nil {
			return players, err
		}

		players[dbPlayer.Name] = domain.Player{
			Team:         dbPlayer.Team,
			InitialOrder: dbPlayer.InitialOrder,
			Order:        dbPlayer.COrder,
			Hand:         hand,
		}
	}

	return players, nil
}
