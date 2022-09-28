package repository

import (
	"coinche/domain"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func getGame(tx *sqlx.Tx, gameID int) (domain.Game, error) {
	var game domain.Game
	var deck []byte

	err := tx.QueryRow(`SELECT * FROM game WHERE id=$1`, gameID).Scan(
		&game.ID,
		&game.Name,
		&game.CreatedAt,
		&game.Phase,
		&deck,
	)

	if err != nil {
		return domain.Game{}, err
	}

	err = json.Unmarshal(deck, &game.Deck)
	if err != nil {
		return domain.Game{}, errors.New(fmt.Sprint(err, "Deck: ", deck))
	}

	game.Players, err = getPlayers(tx, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Bids, err = getBids(tx, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Turns, err = getTurns(tx, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Scores, err = getScoresOrPoints(tx, gameID, "score")
	if err != nil {
		return domain.Game{}, err
	}

	game.Points, err = getScoresOrPoints(tx, gameID, "point")

	return game, err
}

func (s *GameRepository) GetGame(gameID int) (domain.Game, error) {
	tx := s.db.MustBegin()

	game, err := getGame(tx, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	return game, tx.Commit()
}

func (s *GameRepository) ListGames() ([]domain.Game, error) {
	var games []domain.Game
	gamesIds := []int{}

	tx := s.db.MustBegin()

	err := tx.Select(&gamesIds, "SELECT id FROM game")
	if err != nil {
		return nil, err
	}

	for _, gameId := range gamesIds {
		game, err := getGame(tx, gameId)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, tx.Commit()
}
