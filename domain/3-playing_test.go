package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartPlaying(test *testing.T) {
	assert := assert.New(test)

	test.Run("should put the higher bid as trump", func(test *testing.T) {
		game := newBiddingGame()
		game.Deck = []CardID{
			C7, C8, C9, C10, CJ, CQ, CK, CA,
			D7, D8, D9, D10, DJ, DQ, DK, DA,
			H7, H8, H9, H10, HJ, HQ, HK, HA,
			S7, S8, S9, S10, SJ, SQ, SK, SA,
		}

		game.Bids = map[BidValue]Bid{
			Eighty:  {Player: "P1", Color: Club},
			Ninety:  {Player: "P2", Color: Diamond},
			Hundred: {Player: "P3", Color: Heart},
		}

		game.startPlaying()

		assert.Equal(Heart, game.trump())
	})
}

func newPlayingGame() Game {
	return Game{
		ID:   2,
		Name: "GAME TWO",
		Players: map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{C7, C8, C9, DJ, DQ, HJ, HQ, HK}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{C10, CJ, CQ, DK, DA, HA, S7, S8}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{CK, CA, D7, H7, H8, S9, S10, SJ}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []CardID{D8, D9, D10, H9, H10, SQ, SK, SA}},
		},
		Phase: Playing,
		Bids: map[BidValue]Bid{
			Eighty: {Player: "P1", Color: Heart},
		},
	}
}

func TestPlaying(test *testing.T) {
	assert := assert.New(test)

	test.Run("should fail if not in playing", func(test *testing.T) {
		biddingGame := Game{
			ID:    2,
			Phase: Bidding,
		}
		err := biddingGame.Play("P1", C9)

		assert.Error(err)
		assert.Equal(ErrNotPlaying, err.Error())
	})
	test.Run("should fail if not in hand", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C10)

		assert.Error(err)
		assert.Equal(ErrCardNotInHand, err.Error())
	})
	test.Run("should fail if not your turn", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P2", C10)

		assert.Error(err)
		assert.Equal(ErrNotYourTurn+" P2 2", err.Error())
	})
	test.Run("should be able to play a card", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C9)

		assert.NoError(err)
		assert.Equal(7, len(game.Players["P1"].Hand))
		assert.Equal(1, len(game.Turns))

		turn := game.Turns[0]
		assert.Equal(1, len(turn.Plays))
	})
}

func TestEndOfTurn(test *testing.T) {
	assert := assert.New(test)

	test.Run("should end the turn and determin a winner after four plays", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C7)
		assert.NoError(err)

		err = game.Play("P2", CJ)
		assert.NoError(err)

		err = game.Play("P3", CA)
		assert.NoError(err)

		err = game.Play("P4", H9)
		assert.NoError(err)

		assert.Equal(1, len(game.Turns))

		turn := game.Turns[0]
		assert.Equal(4, len(turn.Plays))
		assert.Equal("P4", turn.Winner)
		assert.Equal(1, game.Players["P4"].Order)
	})
}

func TestCanPlay(test *testing.T) {
	assert := assert.New(test)

	test.Run("should fail to play if trying to play another color while having the asked color", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C7)
		assert.NoError(err)

		err = game.Play("P2", S7)

		assert.Error(err)
		assert.Equal(ErrShouldPlayAskedColor, err.Error())
	})

	test.Run("should fail to play if trying to play a trump while having the asked color", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C7)
		assert.NoError(err)

		err = game.Play("P2", HA)

		assert.Error(err)
		assert.Equal(ErrShouldPlayAskedColor, err.Error())
	})

	test.Run("should be able to play anything if has no trump and no asked color", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{C7}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{DK}},
		}

		err := game.Play("P1", C7)
		assert.NoError(err)

		err = game.Play("P2", DK)
		assert.NoError(err)
	})

	test.Run("should fail to play if has trump, no asked color and plays something that is not trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{C7}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{DK, H10}},
		}

		err := game.Play("P1", C7)
		assert.NoError(err)

		err = game.Play("P2", DK)

		assert.Error(err)
		assert.Equal(ErrShouldPlayTrump, err.Error())
	})

	test.Run("should be able to play no trump if partner is winner and trump not asked", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{C10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{C7}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{D7, H8}},
		}

		err := game.Play("P1", C10)
		assert.NoError(err)

		err = game.Play("P2", C7)
		assert.NoError(err)

		err = game.Play("P3", D7)
		assert.NoError(err)
	})

	test.Run("should fail to play no trump if partner is winner and trump is asked", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{HJ}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{C7}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{D7, H8}},
		}

		err := game.Play("P1", HJ)
		assert.NoError(err)

		err = game.Play("P2", C7)
		assert.NoError(err)

		err = game.Play("P3", D7)
		assert.Error(err)
		assert.Equal(ErrShouldPlayAskedColor, err.Error())
	})

	test.Run("should be able to play a bigger trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{C10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{H9}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{HJ, H8}},
		}

		err := game.Play("P1", C10)
		assert.NoError(err)

		err = game.Play("P2", H9)
		assert.NoError(err)

		err = game.Play("P3", HJ)
		assert.NoError(err)
	})

	test.Run("if trump is asked, should be able to play a lower trump if has NO bigger trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{C10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{H9}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{D7, H8}},
		}

		err := game.Play("P1", C10)
		assert.NoError(err)

		err = game.Play("P2", H9)
		assert.NoError(err)

		err = game.Play("P3", H8)
		assert.NoError(err)
	})

	test.Run("if trump is NOT asked, should be able to play a lower trump while having a bigger trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{C10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{H9}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{HJ, H8}},
		}

		err := game.Play("P1", C10)
		assert.NoError(err)

		err = game.Play("P2", H9)
		assert.NoError(err)

		err = game.Play("P3", H8)
		assert.NoError(err)
	})

	test.Run("if trump IS asked, should fail to play a lower trump while having a bigger trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{H9}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{C10}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{HJ, H8}},
		}

		err := game.Play("P1", H9)
		assert.NoError(err)

		err = game.Play("P2", C10)
		assert.NoError(err)

		err = game.Play("P3", H8)
		assert.Error(err)

		assert.Equal(ErrShouldPlayBiggerTrump, err.Error())
	})
}

func TestEndOfPlayingPhase(test *testing.T) {
	assert := assert.New(test)

	test.Run("should pass to counting phase when 8th turn is over", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []CardID{H9}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []CardID{C10}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []CardID{HJ}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []CardID{D8}},
		}
		game.Turns = []Turn{
			{[]Play{}, "P1"},
			{[]Play{}, "P2"},
			{[]Play{}, "P3"},
			{[]Play{}, "P2"},
			{[]Play{}, "P2"},
			{[]Play{}, "P2"},
			{[]Play{
				{"P1", H9},
				{"P2", C10},
				{"P3", HJ},
				{"P4", D8},
			}, "P1"},
		}

		err := game.Play("P1", H9)
		assert.NoError(err)

		err = game.Play("P2", C10)
		assert.NoError(err)

		err = game.Play("P3", HJ)
		assert.NoError(err)

		assert.Equal(Playing, game.Phase)

		err = game.Play("P4", D8)
		assert.NoError(err)

		assert.Equal(Counting, game.Phase)
	})
}
