package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newBiddingGame() Game {
	return Game{
		ID:   2,
		Name: "GAME TWO",
		Players: map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1},
			"P2": {Team: "even", Order: 2, InitialOrder: 2},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3},
			"P4": {Team: "even", Order: 4, InitialOrder: 4},
		},
		Phase: Bidding,
		Bids:  make(map[BidValue]Bid),
		deck:  newDeck(),
	}
}

func TestBidding(test *testing.T) {
	assert := assert.New(test)
	test.Run("should fail if not in bidding", func(test *testing.T) {
		teamingGame := Game{ID: 2, Phase: Teaming}

		err := teamingGame.PlaceBid("P1", Eighty, Spade)

		assert.Error(err)
		assert.Equal(ErrNotBidding, err.Error())
	})

	test.Run("should fail to coinche if no previous bid", func(test *testing.T) {
		game := newBiddingGame()

		err := game.Coinche("P1")

		assert.Error(err)
		assert.Equal(ErrNoBidYet, err.Error())
	})

	test.Run("should place a bid", func(test *testing.T) {
		game := newBiddingGame()

		err := game.PlaceBid("P1", Eighty, Spade)

		want := Bid{Player: "P1", Color: Spade}
		assert.NoError(err)
		assert.Equal(want, game.Bids[Eighty])
	})

	test.Run("should place another bid", func(test *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P1", Color: Spade},
		}
		game.Players = map[string]Player{
			"P1": {Order: 4},
			"P2": {Order: 1},
			"P3": {Order: 2},
			"P4": {Order: 3},
		}

		err := game.PlaceBid("P2", Ninety, Club)

		want := Bid{Player: "P2", Color: Club}
		assert.NoError(err)
		assert.Equal(want, game.Bids[Ninety])
	})

	test.Run("should fail if placing a bid smaller or equal to previous bid", func(test *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P1", Color: Spade},
		}
		game.Players = map[string]Player{
			"P1": {Order: 3},
			"P2": {Order: 4},
			"P3": {Order: 1},
			"P4": {Order: 2},
		}

		err := game.PlaceBid("P3", Eighty, Club)

		assert.Error(err)
		assert.Equal(ErrBidTooSmall, err.Error())
	})

	test.Run("should fail if the player is not the right one", func(test *testing.T) {
		game := newBiddingGame()

		err := game.PlaceBid("P2", Hundred, Club)

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())
	})

	test.Run("order should rotate correctly", func(test *testing.T) {
		game := newBiddingGame()

		err := game.Pass("P1")
		assert.NoError(err)

		err = game.Pass("P2")
		assert.NoError(err)

		err = game.Pass("P3")
		assert.NoError(err)
	})

	test.Run("should fail if the player bid on its own color", func(test *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P1", Color: Spade},
		}

		err := game.PlaceBid("P1", HundredAndTen, Spade)

		assert.Error(err)
		assert.Equal(ErrBiddingItsOwnColor, err.Error())
	})
}

func TestCoinche(test *testing.T) {
	assert := assert.New(test)

	test.Run("same team player should not be able to coinche", func(test *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P1", Color: Spade, Coinche: 1},
		}
		game.Players = map[string]Player{
			"P1": {Order: 4},
			"P2": {Order: 1},
			"P3": {Order: 2},
			"P4": {Order: 3},
		}

		err := game.Coinche("P3")

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())

		err = game.Coinche("P1")

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())
	})

	test.Run("should be able to coinche several times", func(t *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P4", Color: Spade},
		}

		err := game.Coinche("P1")
		assert.NoError(err)

		err = game.Coinche("P4")
		assert.NoError(err)
	})
}

func TestEndOfBidding(test *testing.T) {
	assert := assert.New(test)

	test.Run("should start playing after two passes after coinche", func(t *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P4", Color: Spade, Coinche: 2},
		}

		err := game.Pass("P1")
		assert.NoError(err)

		err = game.Pass("P3")
		assert.NoError(err)

		assert.Equal(Playing, game.Phase)
	})

	test.Run("should start playing after 3 coinches", func(t *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P4", Color: Spade, Coinche: 2},
		}

		err := game.Coinche("P1")

		assert.NoError(err)
		assert.Equal(Playing, game.Phase)
	})

	test.Run("should start playing after 4 passes", func(t *testing.T) {
		game := newBiddingGame()
		game.Bids = map[BidValue]Bid{
			Eighty: {Player: "P4", Color: Spade},
		}

		err := game.Pass("P1")
		assert.NoError(err)
		err = game.Pass("P2")
		assert.NoError(err)
		err = game.Pass("P3")
		assert.NoError(err)
		err = game.Pass("P4")
		assert.NoError(err)

		assert.Equal(Playing, game.Phase)
	})
}
