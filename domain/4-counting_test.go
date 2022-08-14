package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newNormalGame() Game {
	return Game{
		ID:   2,
		Name: "GAME TWO",
		Players: map[string]Player{
			"P1": {Team: "odd", Order: 2, InitialOrder: 1, Hand: []CardID{}},
			"P2": {Team: "even", Order: 3, InitialOrder: 2, Hand: []CardID{}},
			"P3": {Team: "odd", Order: 4, InitialOrder: 3, Hand: []CardID{}},
			"P4": {Team: "even", Order: 1, InitialOrder: 4, Hand: []CardID{S_A}},
		},
		Phase: Counting,
		Bids: map[BidValue]Bid{
			Eighty: {
				Player:  "P1",
				Color:   Heart,
				Coinche: 0,
				Pass:    0,
			},
		},
		scores: map[string]int{
			"odd":  0,
			"even": 0,
		},
		turns: []turn{
			{[]play{
				{"P1", C_7},
				{"P2", C_10},
				{"P3", C_K},
				{"P4", H_9},
			}, "P4"},
			{[]play{
				{"P4", D_8},
				{"P1", D_J},
				{"P2", D_K},
				{"P3", D_7},
			}, "P2"},
			{[]play{
				{"P2", C_J},
				{"P3", C_A},
				{"P4", H_10},
				{"P1", C_8},
			}, "P4"},
			{[]play{
				{"P4", D_9},
				{"P1", D_Q},
				{"P2", D_A},
				{"P3", H_7},
			}, "P3"},
			{[]play{
				{"P3", H_8},
				{"P4", D_10},
				{"P1", H_J},
				{"P2", H_A},
			}, "P1"},
			{[]play{
				{"P1", C_9},
				{"P2", C_Q},
				{"P3", S_9},
				{"P4", S_Q},
			}, "P2"},
			{[]play{
				{"P2", S_7},
				{"P3", S_10},
				{"P4", S_K},
				{"P1", H_Q},
			}, "P1"},
			{[]play{
				{"P1", S_8},
				{"P2", H_K},
				{"P3", S_J},
				{"P4", S_A},
			}, "P2"},
		},
	}
}

func newGameWithBelote() Game {
	game := newNormalGame()
	game.turns = game.turns[:len(game.turns)-1]
	game.turns = append(game.turns, turn{[]play{
		{"P1", H_K}, // belote with H_Q
		{"P2", S_8},
		{"P3", S_J},
		{"P4", S_A},
	}, "P1"})
	return game
}

func newGameWithNoTrump() Game {
	game := newNormalGame()
	game.Bids = map[BidValue]Bid{
		Eighty: {
			Player:  "P1",
			Color:   NoTrump,
			Coinche: 0,
			Pass:    0,
		},
	}
	game.turns = []turn{
		{[]play{
			{"P1", C_7},
			{"P2", C_10},
			{"P3", C_K},
			{"P4", H_9},
		}, "P2"},
		{[]play{
			{"P4", D_8},
			{"P1", D_J},
			{"P2", D_K},
			{"P3", D_7},
		}, "P2"},
		{[]play{
			{"P2", C_J},
			{"P3", C_A},
			{"P4", H_10},
			{"P1", C_8},
		}, "P3"},
		{[]play{
			{"P4", D_9},
			{"P1", D_Q},
			{"P2", D_A},
			{"P3", H_7},
		}, "P2"},
		{[]play{
			{"P3", H_8},
			{"P4", D_10},
			{"P1", H_J},
			{"P2", H_A},
		}, "P2"},
		{[]play{
			{"P1", C_9},
			{"P2", C_Q},
			{"P3", S_9},
			{"P4", S_Q},
		}, "P2"},
		{[]play{
			{"P2", S_7},
			{"P3", S_10},
			{"P4", S_K},
			{"P1", H_Q},
		}, "P3"},
		{[]play{
			{"P1", S_8},
			{"P2", H_K},
			{"P3", S_J},
			{"P4", S_A},
		}, "P4"},
	}
	return game
}

func newGameWithAllTrump() Game {
	game := newNormalGame()
	game.Bids = map[BidValue]Bid{
		Eighty: {
			Player:  "P2",
			Color:   AllTrump,
			Coinche: 0,
			Pass:    0,
		},
	}
	game.turns = []turn{
		{[]play{
			{"P1", C_7},
			{"P2", C_10},
			{"P3", C_K}, // first belote with C_Q
			{"P4", H_9},
		}, "P2"},
		{[]play{
			{"P4", D_8},
			{"P1", D_J},
			{"P2", D_K}, // could be a belote but too late
			{"P3", D_7},
		}, "P1"},
		{[]play{
			{"P2", C_J},
			{"P3", C_A},
			{"P4", H_10},
			{"P1", C_8},
		}, "P2"},
		{[]play{
			{"P4", D_9},
			{"P1", D_A},
			{"P2", D_Q}, // could be a rebelote but too late
			{"P3", H_7},
		}, "P4"},
		{[]play{
			{"P3", H_8},
			{"P4", D_10},
			{"P1", H_J},
			{"P2", H_A},
		}, "P2"},
		{[]play{
			{"P1", C_9},
			{"P2", C_Q}, // rebelote
			{"P3", S_9},
			{"P4", S_Q},
		}, "P1"},
		{[]play{
			{"P2", S_7},
			{"P3", S_10},
			{"P4", S_K},
			{"P1", H_Q},
		}, "P3"},
		{[]play{
			{"P1", S_8},
			{"P2", H_K},
			{"P3", S_J},
			{"P4", S_A},
		}, "P3"},
	}
	return game
}

func newGameWithCapotLost() Game {
	game := newNormalGame()
	game.Bids = map[BidValue]Bid{
		Capot: {
			Player:  "P2",
			Color:   Heart,
			Coinche: 0,
			Pass:    0,
		},
	}
	game.turns = []turn{
		{[]play{
			{"P1", C_7},
			{"P2", C_10},
			{"P3", C_K},
			{"P4", H_9},
		}, "P4"},
		{[]play{
			{"P4", D_8},
			{"P1", D_J},
			{"P2", D_K},
			{"P3", D_7},
		}, "P2"},
		{[]play{
			{"P2", C_J},
			{"P3", C_A},
			{"P4", H_10},
			{"P1", C_8},
		}, "P4"},
		{[]play{
			{"P4", D_9},
			{"P1", D_Q},
			{"P2", H_7},
			{"P3", D_A},
		}, "P2"},
		{[]play{
			{"P3", H_8},
			{"P4", D_10},
			{"P1", H_A},
			{"P2", H_J},
		}, "P2"},
		{[]play{
			{"P1", C_9},
			{"P2", C_Q},
			{"P3", S_9},
			{"P4", S_Q},
		}, "P2"},
		{[]play{
			{"P1", S_8},
			{"P2", H_K},
			{"P3", S_J},
			{"P4", S_A},
		}, "P2"},
		{[]play{
			{"P2", S_7},
			{"P3", S_10},
			{"P4", S_K},
			{"P1", H_Q},
		}, "P1"},
	}
	return game
}

func newGameWithCapotWon() Game {
	game := newGameWithCapotLost()
	game.turns = game.turns[:len(game.turns)-1]
	game.turns = append(game.turns, turn{[]play{
		{"P2", S_7},
		{"P3", S_10},
		{"P4", H_Q},
		{"P1", S_K},
	}, "P4"},
	)
	return game
}

func TestCountingPhase(test *testing.T) {
	assert := assert.New(test)

	test.Run("should go to counting phase on last game", func(test *testing.T) {
		game := newNormalGame()
		game.Phase = Playing
		game.turns = game.turns[:len(game.turns)-1]
		game.turns = append(game.turns, turn{[]play{
			{"P1", S_8},
			{"P2", H_K},
			{"P3", S_J},
		}, ""},
		)

		err := game.Play("P4", S_A)
		assert.NoError(err)

		assert.Equal(Counting, game.Phase)
	})

	/*test.Run("should restart a new game on restart", func(test *testing.T) {
		game := newNormalGame()
		assert.Equal(Counting, game.Phase)

		err := game.Restart()
		assert.NoError(err)
	})*/ // TODO: implement restart
}

func TestPlayersCards(test *testing.T) {
	assert := assert.New(test)

	test.Run("should count correctly in a game", func(test *testing.T) {
		game := newNormalGame()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID{H_8, D_10, H_J, H_A, S_7, S_10, S_K, H_Q}, playerCards["P1"])
		assert.Equal([]CardID{D_8, D_J, D_K, D_7, C_9, C_Q, S_9, S_Q, S_8, H_K, S_J, S_A}, playerCards["P2"])
		assert.Equal([]CardID{D_9, D_Q, D_A, H_7}, playerCards["P3"])
		assert.Equal([]CardID{C_7, C_10, C_K, H_9, C_J, C_A, H_10, C_8}, playerCards["P4"])
	})

	test.Run("should count correctly in a game with BELOTE", func(test *testing.T) {
		game := newGameWithBelote()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID{H_8, D_10, H_J, H_A, S_7, S_10, S_K, H_Q, H_K, S_8, S_J, S_A}, playerCards["P1"])
		assert.Equal([]CardID{D_8, D_J, D_K, D_7, C_9, C_Q, S_9, S_Q}, playerCards["P2"])
		assert.Equal([]CardID{D_9, D_Q, D_A, H_7}, playerCards["P3"])
		assert.Equal([]CardID{C_7, C_10, C_K, H_9, C_J, C_A, H_10, C_8}, playerCards["P4"])
	})

	test.Run("should count correctly in a game with NO-TRUMP", func(test *testing.T) {
		game := newGameWithNoTrump()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID(nil), playerCards["P1"])
		assert.Equal([]CardID{C_7, C_10, C_K, H_9, D_8, D_J, D_K, D_7, D_9, D_Q, D_A, H_7, H_8, D_10, H_J, H_A, C_9, C_Q, S_9, S_Q}, playerCards["P2"])
		assert.Equal([]CardID{C_J, C_A, H_10, C_8, S_7, S_10, S_K, H_Q}, playerCards["P3"])
		assert.Equal([]CardID{S_8, H_K, S_J, S_A}, playerCards["P4"])
	})

	test.Run("should count correctly in a game with ALL-TRUMP and BELOTE (for odd team)", func(test *testing.T) {
		game := newGameWithAllTrump()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID{D_8, D_J, D_K, D_7, C_9, C_Q, S_9, S_Q}, playerCards["P1"])
		assert.Equal([]CardID{C_7, C_10, C_K, H_9, C_J, C_A, H_10, C_8, H_8, D_10, H_J, H_A}, playerCards["P2"])
		assert.Equal([]CardID{S_7, S_10, S_K, H_Q, S_8, H_K, S_J, S_A}, playerCards["P3"])
		assert.Equal([]CardID{D_9, D_A, D_Q, H_7}, playerCards["P4"])
	})
}

func TestCounting(test *testing.T) {
	assert := assert.New(test)

	test.Run("should count correctly in a LOST game", func(test *testing.T) {
		game := newNormalGame()

		game.calculatesTeamPointsAndScores()
		assert.Equal(72, game.points["odd"])
		assert.Equal(90, game.points["even"])
		assert.Equal(72, game.scores["odd"])
		assert.Equal(160+80, game.scores["even"])

		assert.Equal(162, game.points["odd"]+game.points["even"])
	})

	test.Run("should count correctly in a WON game with COINCHE and BELOTE (for odd team)", func(test *testing.T) {
		game := newGameWithBelote()
		game.Bids = map[BidValue]Bid{
			Eighty: {
				Player:  "P1",
				Color:   Heart,
				Coinche: 1,
				Pass:    0,
			},
		}

		game.calculatesTeamPointsAndScores()
		assert.Equal(99+20, game.points["odd"])
		assert.Equal(63, game.points["even"])
		assert.Equal((80+160+20)*2, game.scores["odd"])
		assert.Equal(0, game.scores["even"])

		assert.Equal(182, game.points["odd"]+game.points["even"])
	})

	test.Run("should count correctly in a WON game with BELOTE", func(test *testing.T) {
		game := newGameWithBelote()

		game.calculatesTeamPointsAndScores()
		assert.Equal(99+20, game.points["odd"])
		assert.Equal(63, game.points["even"])
		assert.Equal((80 + 99 + 20), game.scores["odd"])
		assert.Equal(63, game.scores["even"])

		assert.Equal(182, game.points["odd"]+game.points["even"])
	})

	test.Run("should count correctly in a LOST game with NO-TRUMP", func(test *testing.T) {
		game := newGameWithNoTrump()

		game.calculatesTeamPointsAndScores()
		assert.Equal(49, game.points["odd"])
		assert.Equal(113, game.points["even"])
		assert.Equal(49, game.scores["odd"])
		assert.Equal(160+80, game.scores["even"])

		assert.Equal(162, game.points["odd"]+game.points["even"])
	})

	test.Run("should count correctly in a WON game with ALL-TRUMP and BELOTE (for odd team)", func(test *testing.T) {
		game := newGameWithAllTrump()

		game.calculatesTeamPointsAndScores()
		assert.Equal(76, game.points["odd"])
		assert.Equal(86, game.points["even"])
		assert.Equal(20+76, game.scores["odd"]) // the belote was taken by the odd team but the even team had taken the bid
		assert.Equal(80+86, game.scores["even"])

		assert.Equal(162, game.points["odd"]+game.points["even"])
	})

	test.Run("should count correctly in a game with ALL-TRUMP with SURCOINCHE", func(test *testing.T) {
		game := newGameWithAllTrump()
		game.Bids = map[BidValue]Bid{
			Eighty: {
				Player:  "P2",
				Color:   AllTrump,
				Coinche: 2,
				Pass:    0,
			},
		}

		game.calculatesTeamPointsAndScores()
		assert.Equal(76, game.points["odd"])
		assert.Equal(86, game.points["even"])
		assert.Equal(20*4, game.scores["odd"]) // the belote was taken by the odd team but the even team had taken the bid
		assert.Equal((160+80)*4, game.scores["even"])
	})

	test.Run("should count correctly with CAPOT LOST", func(test *testing.T) {
		game := newGameWithCapotLost()

		game.calculatesTeamPointsAndScores()
		assert.Equal(27, game.points["odd"])
		assert.Equal(135, game.points["even"])
		assert.Equal(320, game.scores["odd"])
		assert.Equal(0, game.scores["even"])
	})

	test.Run("should count correctly with CAPOT WON", func(test *testing.T) {
		game := newGameWithCapotWon()

		game.calculatesTeamPointsAndScores()
		assert.Equal(0, game.points["odd"])
		assert.Equal(162, game.points["even"])
		assert.Equal(0, game.scores["odd"])
		assert.Equal(500, game.scores["even"])
	})
}
