package repository

import (
	"coinche/domain"
	"encoding/json"
	"errors"
	"fmt"
)

func (s *GameRepository) GetGame(gameID int) (domain.Game, error) {
	tx := s.db.MustBegin()

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

	rows, err := tx.Query(`SELECT name, team, initialOrder, cOrder, hand FROM player WHERE gameid=$1`, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Players = map[string]domain.Player{}

	for rows.Next() {
		var player domain.Player
		var playerName string
		var hand []byte
		err := rows.Scan(&playerName, &player.Team, &player.InitialOrder, &player.Order, &hand)
		if err != nil {
			return domain.Game{}, err
		}

		err = json.Unmarshal(hand, &player.Hand)
		if err != nil {
			return domain.Game{}, errors.New(fmt.Sprint(err, "Hand: ", hand, "Player: ", playerName))
		}

		game.Players[playerName] = player
	}

	rows, err = tx.Query(`SELECT value, coinche, color, pass, player FROM bid WHERE gameid=$1`, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Bids = map[domain.BidValue]domain.Bid{}

	for rows.Next() {
		var bid domain.Bid
		var bidValue domain.BidValue

		err := rows.Scan(&bidValue, &bid.Coinche, &bid.Color, &bid.Pass, &bid.Player)
		if err != nil {
			return domain.Game{}, err
		}

		game.Bids[bidValue] = bid
	}

	rows, err = tx.Query(`SELECT winner, plays FROM turn WHERE gameid=$1 ORDER BY position`, gameID)
	if err != nil {
		return domain.Game{}, err
	}
	game.Turns = []domain.Turn{}

	for rows.Next() {
		var turn domain.Turn
		var plays []byte

		err := rows.Scan(&turn.Winner, &plays)
		if err != nil {
			return domain.Game{}, err
		}

		err = json.Unmarshal(plays, &turn.Plays)
		if err != nil {
			return domain.Game{}, errors.New(fmt.Sprint(err, "Plays: ", plays))
		}
		game.Turns = append(game.Turns, turn)
	}

	rows, err = tx.Query(`SELECT value, team FROM point WHERE gameid=$1`, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Points = map[string]int{}

	for rows.Next() {
		var value int
		var team string

		err := rows.Scan(&value, &team)
		if err != nil {
			return domain.Game{}, err
		}

		game.Points[team] = value
	}

	rows, err = tx.Query(`SELECT value, team FROM score WHERE gameid=$1`, gameID)
	if err != nil {
		return domain.Game{}, err
	}

	game.Scores = map[string]int{}

	for rows.Next() {
		var value int
		var team string

		err := rows.Scan(&value, &team)
		if err != nil {
			return domain.Game{}, err
		}

		game.Scores[team] = value
	}
	return game, tx.Commit()

}
