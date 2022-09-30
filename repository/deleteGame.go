package repository

import "errors"

const (
	ErrCannotDeleteWithPlayers = "CANNOT DELETE GAME WITH PLAYERS"
)

func (s *GameRepository) DeleteGame(gameID int) error {
	tx := s.db.MustBegin()

	game, err := s.GetGame(gameID)
	if err != nil {
		return err
	}

	playersCount := len(game.Players)
	if playersCount > 0 {
		return errors.New(ErrCannotDeleteWithPlayers)
	}

	_, err = tx.Exec(`DELETE FROM game WHERE id=$1`, gameID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
