package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) UpdatePlayers(id int, players []string, phase domain.Phase) error {
	tx := s.db.MustBegin()

	var currentPlayers []string

	rows, err := tx.Query(`SELECT name FROM player WHERE gameid=$1`, id)
	if err != nil {
		return err
	}

	for rows.Next() {
		var playerName string
		err := rows.Scan(&playerName)
		if err != nil {
			return err
		}

		currentPlayers = append(currentPlayers, playerName)
	}

	for _, currentPlayerName := range currentPlayers {
		shouldDelete := true
		for _, playerName := range players {
			if currentPlayerName == playerName {
				shouldDelete = false
			}
		}

		if shouldDelete {
			_, err := s.db.Exec(
				`
				DELETE FROM player
				WHERE name = $1 AND gameid = $2
				`,
				currentPlayerName,
				id,
			)
			if err != nil {
				return err
			}
		}
	}

	for _, playerName := range players {
		shouldCreate := true
		for _, currentPlayerName := range currentPlayers {
			if currentPlayerName == playerName {
				shouldCreate = false
			}
		}

		if shouldCreate {
			_, err := s.db.Exec(
				`
				INSERT INTO player (name, gameid) 
				VALUES ($1, $2)
				`,
				playerName,
				id,
			)
			if err != nil {
				return err
			}
		}
	}

	_, err = tx.Exec(
		`
		UPDATE game
		SET phase = $1
		WHERE id = $2
		`,
		phase,
		id,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}
