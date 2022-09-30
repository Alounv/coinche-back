package repository

import (
	"coinche/domain"

	"github.com/jmoiron/sqlx"
)

var bidSchema = `
CREATE TABLE IF NOT EXISTS bid (
	id serial PRIMARY KEY NOT NULL,
	gameid integer NOT NULL REFERENCES game(id),
	value integer NOT NULL,
	player  text NOT NULL,
	color  text NOT NULL,
	coinche integer DEFAULT 0,
	pass integer DEFAULT 0
)`

func createBids(tx *sqlx.Tx, gameID int, bids map[domain.BidValue]domain.Bid) error {
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

func getBids(tx *sqlx.Tx, gameID int) (map[domain.BidValue]domain.Bid, error) {
	type dbBid struct {
		Value   domain.BidValue
		Player  string
		Coinche int
		Color   domain.Color
		Pass    int
	}

	var dbBids []dbBid

	err := tx.Select(&dbBids, `
    SELECT value, player, coinche, color, pass FROM bid WHERE gameid = $1
  `, gameID)
	if err != nil {
		return nil, err
	}

	bids := map[domain.BidValue]domain.Bid{}
	for _, bid := range dbBids {
		bids[bid.Value] = domain.Bid{
			Player:  bid.Player,
			Coinche: bid.Coinche,
			Color:   bid.Color,
			Pass:    bid.Pass,
		}
	}

	return bids, nil
}
