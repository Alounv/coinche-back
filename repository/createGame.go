package repository

import (
	"coinche/domain"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func createGame(game domain.Game, tx *sqlx.Tx) (int, error) {
	var gameID int

	deck, err := json.Marshal(game.Deck)
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow(
		`
		INSERT INTO game (name, phase, deck, root) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
		`,
		game.Name,
		game.Phase,
		deck,
		game.ID,
	).Scan(&gameID)
	if err != nil {
		return 0, err
	}

	for playerName, player := range game.Players {
		err := createFullPlayerForTest(tx, gameID, playerName, player)
		if err != nil {
			return 0, err
		}
	}

	err = updateBids(tx, gameID, game.Bids)
	if err != nil {
		return 0, err
	}

	err = createTurns(tx, gameID, game.Turns)
	if err != nil {
		return 0, err
	}

	err = createAndUpdateScores(tx, gameID, game.Scores)
	if err != nil {
		return 0, err
	}

	err = createPoints(tx, gameID, game.Points)
	if err != nil {
		return 0, err
	}

	return gameID, nil
}

func (s *GameRepository) CreateGame(game domain.Game) (int, error) {
	tx := s.db.MustBegin()

	gameID, err := createGame(game, tx)
	if err != nil {
		return 0, err
	}

	return gameID, tx.Commit()
}
