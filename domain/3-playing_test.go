package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartPlaying(test *testing.T) {
	assert := assert.New(test)

	test.Run("should distribute as expected", func(test *testing.T) {
		game := newBiddingGame()
		game.deck = []cardID{
			C_7, C_8, C_9, C_10, C_J, C_Q, C_K, C_A,
			D_7, D_8, D_9, D_10, D_J, D_Q, D_K, D_A,
			H_7, H_8, H_9, H_10, H_J, H_Q, H_K, H_A,
			S_7, S_8, S_9, S_10, S_J, S_Q, S_K, S_A,
		}

		game.startPlaying()

		assert.Equal([]cardID{}, game.deck)
		assert.Equal([]cardID{C_7, C_8, C_9, D_J, D_Q, H_J, H_Q, H_K}, game.Players["P1"].Hand)
		assert.Equal([]cardID{C_10, C_J, C_Q, D_K, D_A, H_A, S_7, S_8}, game.Players["P2"].Hand)
		assert.Equal([]cardID{C_K, C_A, D_7, H_7, H_8, S_9, S_10, S_J}, game.Players["P3"].Hand)
		assert.Equal([]cardID{D_8, D_9, D_10, H_9, H_10, S_Q, S_K, S_A}, game.Players["P4"].Hand)
	})

	test.Run("should put the higher bid as trump", func(test *testing.T) {
		game := newBiddingGame()
		game.deck = []cardID{
			C_7, C_8, C_9, C_10, C_J, C_Q, C_K, C_A,
			D_7, D_8, D_9, D_10, D_J, D_Q, D_K, D_A,
			H_7, H_8, H_9, H_10, H_J, H_Q, H_K, H_A,
			S_7, S_8, S_9, S_10, S_J, S_Q, S_K, S_A,
		}

		game.Bids = map[BidValue]Bid{
			Eighty:  {Player: "P1", Color: Club},
			Ninety:  {Player: "P2", Color: Diamond},
			Hundred: {Player: "P3", Color: Heart},
		}

		game.startPlaying()

		assert.Equal(Heart, game.trump)
	})
}

func newPlayingGame() Game {
	return Game{
		ID:   2,
		Name: "GAME TWO",
		Players: map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{C_7, C_8, C_9, D_J, D_Q, H_J, H_Q, H_K}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{C_10, C_J, C_Q, D_K, D_A, H_A, S_7, S_8}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{C_K, C_A, D_7, H_7, H_8, S_9, S_10, S_J}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8, D_9, D_10, H_9, H_10, S_Q, S_K, S_A}},
		},
		Phase: Playing,
		Bids:  map[BidValue]Bid{},
		trump: Heart,
	}
}

func TestPlaying(test *testing.T) {
	assert := assert.New(test)

	test.Run("should fail if not in playing", func(test *testing.T) {
		biddingGame := Game{
			ID:    2,
			Phase: Bidding,
		}
		err := biddingGame.Play("P1", C_9)

		assert.Error(err)
		assert.Equal(ErrNotPlaying, err.Error())
	})
	test.Run("should fail if not in hand", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C_10)

		assert.Error(err)
		assert.Equal(ErrCardNotInHand, err.Error())
	})
	test.Run("should fail if not your turn", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P2", C_10)

		assert.Error(err)
		assert.Equal(ErrNotYourTurn, err.Error())
	})
	test.Run("should be able to play a card", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C_9)

		assert.NoError(err)
		assert.Equal(7, len(game.Players["P1"].Hand))
		assert.Equal(1, len(game.turns))

		turn := game.turns[0]
		assert.Equal(1, len(turn.plays))
	})
}

func TestEndOfTurn(test *testing.T) {
	assert := assert.New(test)

	test.Run("should end the turn and determin a winner after four plays", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C_7)
		assert.NoError(err)

		err = game.Play("P2", C_J)
		assert.NoError(err)

		err = game.Play("P3", C_A)
		assert.NoError(err)

		err = game.Play("P4", H_9)
		assert.NoError(err)

		assert.Equal(1, len(game.turns))

		turn := game.turns[0]
		assert.Equal(4, len(turn.plays))
		assert.Equal("P4", turn.winner)
		assert.Equal(1, game.Players["P4"].Order)
	})
}

func TestCanPlay(test *testing.T) {
	assert := assert.New(test)

	test.Run("should fail to play if trying to play another color while having the asked color", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C_7)
		assert.NoError(err)

		err = game.Play("P2", S_7)

		assert.Error(err)
		assert.Equal(ErrShouldPlayAskedColor, err.Error())
	})

	test.Run("should fail to play if trying to play a trump while having the asked color", func(test *testing.T) {
		game := newPlayingGame()

		err := game.Play("P1", C_7)
		assert.NoError(err)

		err = game.Play("P2", H_A)

		assert.Error(err)
		assert.Equal(ErrShouldPlayAskedColor, err.Error())
	})

	test.Run("should be able to play anything if has no trump and no asked color", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{C_7}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{D_K}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{C_K}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.trump = Heart

		err := game.Play("P1", C_7)
		assert.NoError(err)

		err = game.Play("P2", D_K)
		assert.NoError(err)
	})

	test.Run("should fail to play if has trump, no asked color and plays something that is not trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{C_7}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{D_K, H_10}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{C_K}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.trump = Heart

		err := game.Play("P1", C_7)
		assert.NoError(err)

		err = game.Play("P2", D_K)

		assert.Error(err)
		assert.Equal(ErrShouldPlayTrump, err.Error())
	})

	test.Run("should be able to play no trump if partner is winner", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{C_10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{C_7}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{D_7, H_8}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.trump = Heart

		err := game.Play("P1", C_10)
		assert.NoError(err)

		err = game.Play("P2", C_7)
		assert.NoError(err)

		err = game.Play("P3", D_7)
		assert.NoError(err)
	})

	test.Run("should be able to play a bigger trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{C_10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{H_9}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{H_J, H_8}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.trump = Heart

		err := game.Play("P1", C_10)
		assert.NoError(err)

		err = game.Play("P2", H_9)
		assert.NoError(err)

		err = game.Play("P3", H_J)
		assert.NoError(err)
	})

	test.Run("should be able to play a lower trump if has not bigger trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{C_10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{H_9}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{D_7, H_8}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.trump = Heart

		err := game.Play("P1", C_10)
		assert.NoError(err)

		err = game.Play("P2", H_9)
		assert.NoError(err)

		err = game.Play("P3", H_8)
		assert.NoError(err)
	})

	test.Run("should fail to play a lower trump if has a bigger trump", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{C_10}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{H_9}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{H_J, H_8}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.trump = Heart

		err := game.Play("P1", C_10)
		assert.NoError(err)

		err = game.Play("P2", H_9)
		assert.NoError(err)

		err = game.Play("P3", H_8)

		assert.Error(err)
		assert.Equal(ErrShouldPlayBiggerTrump, err.Error())
	})

	test.Run("should be able to play a lower while having a bigger trump if partner is winner", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{H_9}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{C_10}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{H_J, H_8}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.trump = Heart

		err := game.Play("P1", H_9)
		assert.NoError(err)

		err = game.Play("P2", C_10)
		assert.NoError(err)

		err = game.Play("P3", H_8)
		assert.NoError(err)
	})
}

func TestEndOfPlayingPhase(test *testing.T) {
	assert := assert.New(test)

	test.Run("should pass to counting phase when 8th turn is over", func(test *testing.T) {
		game := newPlayingGame()
		game.Players = map[string]Player{
			"P1": {Team: "odd", Order: 1, InitialOrder: 1, Hand: []cardID{H_9}},
			"P2": {Team: "even", Order: 2, InitialOrder: 2, Hand: []cardID{C_10}},
			"P3": {Team: "odd", Order: 3, InitialOrder: 3, Hand: []cardID{H_J}},
			"P4": {Team: "even", Order: 4, InitialOrder: 4, Hand: []cardID{D_8}},
		}
		game.turns = []turn{
			{[]play{}, "P1"},
			{[]play{}, "P2"},
			{[]play{}, "P3"},
			{[]play{}, "P2"},
			{[]play{}, "P2"},
			{[]play{}, "P2"},
			{[]play{
				{"P1", H_9},
				{"P2", C_10},
				{"P3", H_J},
				{"P4", D_8},
			}, "P1"},
		}

		err := game.Play("P1", H_9)
		assert.NoError(err)

		err = game.Play("P2", C_10)
		assert.NoError(err)

		err = game.Play("P3", H_J)
		assert.NoError(err)

		assert.Equal(Playing, game.Phase)

		err = game.Play("P4", D_8)
		assert.NoError(err)

		assert.Equal(Counting, game.Phase)
	})
}
