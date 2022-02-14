package gamerepo

import (
	"coinche/domain"
)

func (s *GameRepository) UpdatePlayers(id int, players map[string]domain.Player, phase domain.Phase) error {
	tx := s.db.MustBegin()

	var currentPlayers map[string]int = map[string]int{}

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

		currentPlayers[playerName] = 0
	}

	for currentPlayerName := range currentPlayers {
		shouldDelete := false
		if _, ok := players[currentPlayerName]; !ok {
			shouldDelete = true
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

	for playerName, player := range players {
		shouldCreate := false
		if _, ok := currentPlayers[playerName]; !ok {
			shouldCreate = true
		}

		if shouldCreate {
			_, err := s.db.Exec(
				`
				INSERT INTO player (name, team, gameid) 
				VALUES ($1, $2, $3)
				`,
				playerName,
				player.Team,
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
