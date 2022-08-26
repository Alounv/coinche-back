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

	return game, tx.Commit()
}
