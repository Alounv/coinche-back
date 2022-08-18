package repository

import (
	"coinche/domain"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func getGame(rows *sql.Rows, db *sqlx.DB) (domain.Game, error) {
	var game domain.Game
	var deck []byte
	err := rows.Scan(
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

	rows, err = db.Query("SELECT name, team FROM player WHERE gameid=$1", game.ID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Players = map[string]domain.Player{}

	for rows.Next() {
		var playerName string
		var teamName string
		err = rows.Scan(&playerName, &teamName)
		game.Players[playerName] = domain.Player{Team: teamName}
		if err != nil {
			return domain.Game{}, err
		}
	}

	return game, nil
}

func (s *GameRepository) ListGames() ([]domain.Game, error) {
	var games []domain.Game

	rows, err := s.db.Query("SELECT * FROM game ")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		game, err := getGame(rows, s.db)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}
