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

func updateBid(tx *sqlx.Tx, gameID int, bidValue domain.BidValue, bid domain.Bid) error {
	_, err := tx.Exec(
		`
    UPDATE bid
    SET  player = $3, coinche = $4, color = $5, pass = $6
    WHERE gameid = $1 AND value = $2
    `,
		gameID,
		bidValue,
		bid.Player,
		bid.Coinche,
		bid.Color,
		bid.Pass,
	)
	return err
}

func createBid(tx *sqlx.Tx, gameID int, bidValue domain.BidValue, bid domain.Bid) error {
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
	return err
}

func updateBids(tx *sqlx.Tx, gameID int, bids map[domain.BidValue]domain.Bid) error {
	currentBids, err := getBids(tx, gameID)
	if err != nil {
		return err
	}
	for bidValue, bid := range bids {
		shouldCreate := false
		if _, ok := currentBids[bidValue]; !ok {
			shouldCreate = true
		}

		if shouldCreate {
			err := createBid(tx, gameID, bidValue, bid)
			if err != nil {
				return err
			}
		} else {
			err := updateBid(tx, gameID, bidValue, bid)
			if err != nil {
				return err
			}
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
