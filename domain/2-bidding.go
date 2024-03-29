package domain

import (
	"errors"
	"sort"
)

const (
	ErrNotBidding         = "NOT IN BIDDING PHASE"
	ErrBidTooSmall        = "BID IS TOO SMALL"
	ErrNotYourTurn        = "NOT YOUR TURN"
	ErrNotYourTeamTurn    = "NOT YOUR TEAM TURN"
	ErrHasBeenCoinched    = "HAS BEEN COINCHED"
	ErrBiddingItsOwnColor = "BIDDING ITS OWN COLOR"
	ErrNoBidYet           = "NO BID YET"
)

func (game *Game) startBidding() error {
	err := game.canStartBidding()
	if err != nil {
		return err
	}
	game.Phase = Bidding

	shouldInitiateOrder := false

	for _, player := range game.Players {
		if player.InitialOrder == 0 {
			shouldInitiateOrder = true
		}
		break
	}

	if shouldInitiateOrder {
		game.initiateOrder()
	} else {
		game.rotateInitialOrder()
	}

	game.distributeCards()

	return nil
}

func (game *Game) distributeCards() {
	for name, player := range game.Players {
		player.Hand = game.draw(player.Order)
		game.Players[name] = player
	}
	game.Deck = []CardID{}
}

func (game *Game) draw(order int) []CardID {
	base := order - 1

	deckIndexes := []int{
		base*3 + 0,
		base*3 + 1,
		base*3 + 2,

		12 + base*2 + 0,
		12 + base*2 + 1,

		20 + base*3 + 0,
		20 + base*3 + 1,
		20 + base*3 + 2,
	}
	hand := []CardID{}
	for _, deckIndex := range deckIndexes {
		hand = append(hand, game.Deck[deckIndex])
	}

	return hand
}

func (game *Game) rotateInitialOrder() {
	for name, player := range game.Players {
		if player.InitialOrder == 1 {
			player.InitialOrder = 4
		} else {
			player.InitialOrder--
		}

		game.Players[name] = player
	}

	game.resetOrderAsInitialOrder()
}

func (game *Game) resetOrderAsInitialOrder() {
	for name, player := range game.Players {
		player.Order = player.InitialOrder
		game.Players[name] = player
	}
}

func (game *Game) initiateOrder() {
	team1 := ""
	team2 := ""

	keys := make([]string, 0, 4)
	for name := range game.Players {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		player := game.Players[name]

		if team1 == "" {
			team1 = player.Team
			player.Order = 1
			player.InitialOrder = 1
		} else if team1 == player.Team {
			player.Order = 3
			player.InitialOrder = 3
		} else if team2 == "" {
			team2 = player.Team
			player.Order = 2
			player.InitialOrder = 2
		} else {
			player.Order = 4
			player.InitialOrder = 4
		}

		game.Players[name] = player
	}
}

func (game *Game) getLastBid() (Bid, BidValue) {
	var maxValue BidValue
	for value := range game.Bids {
		if value > maxValue {
			maxValue = value
		}
	}
	return game.Bids[maxValue], maxValue
}

func (game *Game) PlaceBid(player string, value BidValue, color Color) error {
	if game.Phase != Bidding {
		return errors.New(ErrNotBidding)
	}
	lastBid, maxValue := game.getLastBid()

	if value <= maxValue {
		return errors.New(ErrBidTooSmall)
	}

	err := game.checkPlayerTurn(player)
	if err != nil {
		return err
	}

	if lastBid.Coinche > 0 {
		return errors.New(ErrHasBeenCoinched)
	}

	if lastBid.Player == player && lastBid.Color == color {
		return errors.New(ErrBiddingItsOwnColor)
	}

	game.Bids[value] = Bid{
		Player:  player,
		Color:   color,
		Coinche: 0,
	}

	game.rotateOrder()
	return nil
}

func (game *Game) Pass(player string) error {
	if game.Phase != Bidding {
		return errors.New(ErrNotBidding)
	}

	lastBid, maxValue := game.getLastBid()

	if lastBid.Coinche > 0 && lastBid.Pass == 0 { // In this case any player of the team can pass
		err := game.checkTeamTurn(player)
		if err != nil {
			return err
		}
	} else {
		err := game.checkPlayerTurn(player)
		if err != nil {
			return err
		}
	}

	game.Bids[maxValue] = Bid{
		Player:  lastBid.Player,
		Color:   lastBid.Color,
		Coinche: lastBid.Coinche,
		Pass:    lastBid.Pass + 1,
	}

	if lastBid.Coinche > 0 {
		if lastBid.Pass+1 > 1 {
			game.startPlaying()
			return nil
		}

		if game.Players[player].Order == 1 {
			game.rotateOrder()
			game.rotateOrder() // If the player to pass was the correct one (1) we rotate twice so the second player of the same team can play.
		}
		// If the player was other one (3) we do not rotate because the other player of the same team is already 1.

		return nil
	}

	if lastBid.Pass+1 > 3 {
		game.startPlaying()
		return nil
	}

	game.rotateOrder()
	return nil
}

func (game *Game) Coinche(player string) error {
	if game.Phase != Bidding {
		return errors.New(ErrNotBidding)
	}

	err := game.checkTeamTurn(player)
	if err != nil {
		return err
	}

	lastBid, maxValue := game.getLastBid()

	if maxValue == 0 {
		return errors.New(ErrNoBidYet)
	}

	game.Bids[maxValue] = Bid{
		Player:  lastBid.Player,
		Color:   lastBid.Color,
		Coinche: lastBid.Coinche + 1,
		Pass:    0,
	}

	if lastBid.Coinche+1 > 2 {
		game.startPlaying()
	}

	game.rotateOrder()

	return nil
}
