package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetting(test *testing.T) {
	assert := assert.New(test)
	testGame := Game{
		ID:   2,
		Name: "GAME TWO",
		Players: map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1},
			"P2": {Team: "even", Order: 2, InitialOrder: 2},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3},
			"P4": {Team: "even", Order: 4, InitialOrder: 4},
		},
		Phase: Bidding,
		Bids:  map[BidValue]Bid{},
	}
	test.Run("should fail if not in bidding", func(test *testing.T) {
		teamingGame := Game{
			ID:    2,
			Phase: Teaming,
		}
		err := teamingGame.PlaceBid("P1", Eighty, Spade)

		assert.Error(err)
		assert.Equal(ErrNotBidding, err.Error())
	})

	test.Run("should fail to coinche if no previous bid", func(test *testing.T) {
		err := testGame.Coinche("P1")

		assert.Error(err)
		assert.Equal(ErrNoBidYet, err.Error())
	})

	test.Run("should place a bid", func(test *testing.T) {
		want := Bid{Player: "P1", Color: Spade}
		err := testGame.PlaceBid("P1", Eighty, Spade)

		assert.NoError(err)
		assert.Equal(want, testGame.Bids[Eighty])
	})

	test.Run("should place another bid", func(test *testing.T) {
		want := Bid{Player: "P2", Color: Club}
		err := testGame.PlaceBid("P2", Ninety, Club)

		assert.NoError(err)
		assert.Equal(want, testGame.Bids[Ninety])
	})

	test.Run("should fail if placing a bid smaller or equal to previous bid", func(test *testing.T) {
		err := testGame.PlaceBid("P3", Eighty, Club)

		assert.Error(err)
		assert.Equal(ErrBidTooSmall, err.Error())
	})

	test.Run("should fail if the player is not the right one", func(test *testing.T) {
		err := testGame.PlaceBid("P1", Hundred, Club)

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())
	})

	test.Run("order should rotate correctly", func(test *testing.T) {
		err := testGame.Pass("P3")
		assert.NoError(err)

		err = testGame.Pass("P4")
		assert.NoError(err)

		err = testGame.Pass("P1")
		assert.NoError(err)
	})

	test.Run("should fail if the player bid on its own color", func(test *testing.T) {
		err := testGame.PlaceBid("P2", HundredAndTen, Club)

		assert.Error(err)
		assert.Equal(ErrBiddingItsOwnColor, err.Error())
	})

	test.Run("same team player should not be able to coinche", func(test *testing.T) {
		err := testGame.PlaceBid("P2", HundredAndTen, Heart)
		if err != nil {
			test.Fatal(err)
		}

		err = testGame.Coinche("P2")

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())

		err = testGame.Coinche("P4")

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())
	})

	test.Run("should be able to coinche several times", func(t *testing.T) {
		testGame := Game{
			ID:   2,
			Name: "GAME TWO",
			Players: map[string]Player{
				"P1": {Team: "odd", Order: 1, InitialOrder: 1},
				"P2": {Team: "even", Order: 2, InitialOrder: 2},
				"P3": {Team: "odd", Order: 3, InitialOrder: 3},
				"P4": {Team: "even", Order: 4, InitialOrder: 4},
			},
			Phase: Bidding,
			Bids: map[BidValue]Bid{
				Eighty: {Player: "P4", Color: Spade},
			},
		}

		err := testGame.Coinche("P1")
		assert.NoError(err)

		err = testGame.Coinche("P4")
		assert.NoError(err)
	})

	test.Run("should start playing after two passes after coinche", func(t *testing.T) {
		testGame := Game{
			ID:   2,
			Name: "GAME TWO",
			Players: map[string]Player{
				"P1": {Team: "odd", Order: 1, InitialOrder: 1},
				"P2": {Team: "even", Order: 2, InitialOrder: 2},
				"P3": {Team: "odd", Order: 3, InitialOrder: 3},
				"P4": {Team: "even", Order: 4, InitialOrder: 4},
			},
			Phase: Bidding,
			Bids: map[BidValue]Bid{
				Eighty: {Player: "P4", Color: Spade, Coinche: 2},
			},
			deck: []int{
				31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17,
				16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
			},
		}

		err := testGame.Pass("P1")
		assert.NoError(err)

		err = testGame.Pass("P3")
		assert.NoError(err)

		assert.Equal(Playing, testGame.Phase)
	})

	test.Run("should start playing after 3 coinches", func(t *testing.T) {
		testGame := Game{
			ID:   2,
			Name: "GAME TWO",
			Players: map[string]Player{
				"P1": {Team: "odd", Order: 1, InitialOrder: 1},
				"P2": {Team: "even", Order: 2, InitialOrder: 2},
				"P3": {Team: "odd", Order: 3, InitialOrder: 3},
				"P4": {Team: "even", Order: 4, InitialOrder: 4},
			},
			Phase: Bidding,
			Bids: map[BidValue]Bid{
				Eighty: {Player: "P4", Color: Spade, Coinche: 2},
			},
			deck: []int{
				31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17,
				16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
			},
		}

		err := testGame.Coinche("P1")

		assert.NoError(err)
		assert.Equal(Playing, testGame.Phase)
	})

	test.Run("should start playing after 4 passes", func(t *testing.T) {
		testGame := Game{
			ID:   2,
			Name: "GAME TWO",
			Players: map[string]Player{
				"P1": {Team: "odd", Order: 1, InitialOrder: 1},
				"P2": {Team: "even", Order: 2, InitialOrder: 2},
				"P3": {Team: "odd", Order: 3, InitialOrder: 3},
				"P4": {Team: "even", Order: 4, InitialOrder: 4},
			},
			Phase: Bidding,
			Bids: map[BidValue]Bid{
				Eighty: {Player: "P4", Color: Spade},
			},
			deck: []int{
				31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17,
				16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
			},
		}

		err := testGame.Pass("P1")
		assert.NoError(err)
		err = testGame.Pass("P2")
		assert.NoError(err)
		err = testGame.Pass("P3")
		assert.NoError(err)
		err = testGame.Pass("P4")
		assert.NoError(err)

		assert.Equal(Playing, testGame.Phase)
	})
}
