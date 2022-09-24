package repository

import (
	"coinche/domain"
	"encoding/json"
	"github.com/jmoiron/sqlx"
)

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
