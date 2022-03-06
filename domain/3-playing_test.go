package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartPlaying(test *testing.T) {
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
		Phase: Playing,
		Bids:  map[BidValue]Bid{},
		deck: []int{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
			17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
		},
	}

	test.Run("should distribute as expected", func(test *testing.T) {
		testGame.startPlaying()

		assert.Equal([]int{}, testGame.deck)
		assert.Equal([]int{0, 1, 2, 12, 13, 20, 21, 22}, testGame.Players["P1"].Hand)
		assert.Equal([]int{3, 4, 5, 14, 15, 23, 24, 25}, testGame.Players["P2"].Hand)
		assert.Equal([]int{6, 7, 8, 16, 17, 26, 27, 28}, testGame.Players["P3"].Hand)
		assert.Equal([]int{9, 10, 11, 18, 19, 29, 30, 31}, testGame.Players["P4"].Hand)
	})

	test.Run("should distribute as expected with disordered values", func(test *testing.T) {
		testGame.deck = []int{
			31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17,
			16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
		}

		testGame.startPlaying()

		assert.Equal([]int{}, testGame.deck)
		assert.Equal([]int{31, 30, 29, 19, 18, 11, 10, 9}, testGame.Players["P1"].Hand)
		assert.Equal([]int{28, 27, 26, 17, 16, 8, 7, 6}, testGame.Players["P2"].Hand)
		assert.Equal([]int{25, 24, 23, 15, 14, 5, 4, 3}, testGame.Players["P3"].Hand)
		assert.Equal([]int{22, 21, 20, 13, 12, 2, 1, 0}, testGame.Players["P4"].Hand)
	})

	test.Run("should put the higher bid as trump", func(test *testing.T) {
		testGame.deck = []int{
			31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17,
			16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
		}

		testGame.Bids = map[BidValue]Bid{
			Eighty:  {Player: "P1", Color: Club},
			Ninety:  {Player: "P2", Color: Diamond},
			Hundred: {Player: "P3", Color: Heart},
		}

		testGame.startPlaying()

		assert.Equal(Heart, testGame.trump)
	})
}

func TestPlaying(test *testing.T) {
	assert := assert.New(test)
	testGame := Game{
		ID:   2,
		Name: "GAME TWO",
		Players: map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []int{0, 1, 2, 12, 13, 20, 21, 22}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []int{3, 4, 5, 14, 15, 23, 24, 25}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []int{6, 7, 8, 16, 17, 26, 27, 28}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []int{9, 10, 11, 18, 19, 29, 30, 31}},
		},
		Phase: Playing,
		Bids:  map[BidValue]Bid{},
	}
	test.Run("should fail if not in playing", func(test *testing.T) {
		biddingGame := Game{
			ID:    2,
			Phase: Bidding,
		}
		err := biddingGame.Play("P1", 0)

		assert.Error(err)
		assert.Equal(ErrNotPlaying, err.Error())
	})
	test.Run("should fail if not in hand", func(test *testing.T) {
		err := testGame.Play("P1", 3)

		assert.Error(err)
		assert.Equal(ErrCardNotInHand, err.Error())
	})
	test.Run("should fail if not your turn", func(test *testing.T) {
		err := testGame.Play("P2", 3)

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())
	})
	test.Run("should be able to play a card", func(test *testing.T) {
		err := testGame.Play("P1", 0)

		assert.NoError(err)
		assert.Equal(7, len(testGame.Players["P1"].Hand))
		assert.Equal(1, len(testGame.turns))

		turn := testGame.turns[0]
		assert.Equal(1, len(turn.plays))
	})
	test.Run("should end the turn and determin a winner after four plays", func(test *testing.T) {
		err := testGame.Play("P2", 3)
		assert.NoError(err)
		err = testGame.Play("P3", 6)
		assert.NoError(err)
		err = testGame.Play("P4", 9)

		assert.NoError(err)
		assert.Equal(1, len(testGame.turns))

		turn := testGame.turns[0]
		assert.Equal(4, len(turn.plays))
		assert.Equal("P2", turn.winner)
		assert.Equal(1, testGame.Players["P2"].Order)
	})
}
