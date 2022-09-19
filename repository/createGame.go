package repository

import (
	"coinche/domain"
	"coinche/utilities"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func (s *GameRepository) CreateBids(tx *sqlx.Tx, gameID int, bids map[domain.BidValue]domain.Bid) error {
	for bidValue, bid := range bids {
		_, err := tx.Exec(
			`
			INSERT INTO bid (gameid, value, player, coinche, color, pass) 
			VALUES ($1, $2, $3, $4, $5, $6)
			`,
			gameID,
			bidValue,
			bid.Player,
			bid.Coinche,
			bid.Color,
			bid.Pass,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GameRepository) CreateGame(game domain.Game) (int, error) {
	tx := s.db.MustBegin()

	gameID, err := s.createAGame(game, tx)
	if err != nil {
		return 0, err
	}

	return gameID, tx.Commit()
}

func (s *GameRepository) createAGame(game domain.Game, tx *sqlx.Tx) (int, error) {
	var gameID int

	deck, err := json.Marshal(game.Deck)
	utilities.PanicIfErr(err)

	err = tx.QueryRow(
		`
		INSERT INTO game (name, phase, deck) 
		VALUES ($1, $2, $3) 
		RETURNING id
		`,
		game.Name,
		game.Phase,
		deck,
	).Scan(&gameID)
	if err != nil {
		return 0, err
	}

	for playerName, player := range game.Players {
		hand, err := json.Marshal(player.Hand)
		utilities.PanicIfErr(err)

		_, err = tx.Exec(
			`
			INSERT INTO player (name, team, gameid, initialOrder, cOrder, hand) 
			VALUES ($1, $2, $3, $4, $5, $6)
			`,
			playerName,
			player.Team,
			gameID,
			player.InitialOrder,
			player.Order,
			hand,
		)
		if err != nil {
			return 0, err
		}
	}

	err = s.CreateBids(tx, gameID, game.Bids)
	if err != nil {
		return 0, err
	}

	for _, turn := range game.Turns {
		var plays []byte
		plays, err = json.Marshal(turn.Plays)
		utilities.PanicIfErr(err)
		_, err = tx.Exec(
			`
			INSERT INTO turn (gameid, winner, plays) 
			VALUES ($1, $2, $3)
			`,
			gameID,
			turn.Winner,
			plays,
		)

		if err != nil {
			return 0, err
		}
	}

	for team, points := range game.Points {
		utilities.PanicIfErr(err)
		_, err = tx.Exec(
			`
			INSERT INTO point (gameid, team, value) 
			VALUES ($1, $2, $3)
			`,
			gameID,
			team,
			points,
		)

		if err != nil {
			return 0, err
		}
	}

	for team, scores := range game.Scores {
		utilities.PanicIfErr(err)
		_, err = tx.Exec(
			`
			INSERT INTO score (gameid, team, value) 
			VALUES ($1, $2, $3)
			`,
			gameID,
			team,
			scores,
		)

		if err != nil {
			return 0, err
		}
	}

	return gameID, nil
}
