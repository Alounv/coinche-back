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
			"P4": {Team: "even", Order: 1, InitialOrder: 4, Hand: []CardID{SA}},
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
		Scores: map[string]int{
			"odd":  0,
			"even": 0,
		},
		Turns: []Turn{
			{[]Play{
				{"P1", C7},
				{"P2", C10},
				{"P3", CK},
				{"P4", H9},
			}, "P4"},
			{[]Play{
				{"P4", D8},
				{"P1", DJ},
				{"P2", DK},
				{"P3", D7},
			}, "P2"},
			{[]Play{
				{"P2", CJ},
				{"P3", CA},
				{"P4", H10},
				{"P1", C8},
			}, "P4"},
			{[]Play{
				{"P4", D9},
				{"P1", DQ},
				{"P2", DA},
				{"P3", H7},
			}, "P3"},
			{[]Play{
				{"P3", H8},
				{"P4", D10},
				{"P1", HJ},
				{"P2", HA},
			}, "P1"},
			{[]Play{
				{"P1", C9},
				{"P2", CQ},
				{"P3", S9},
				{"P4", SQ},
			}, "P2"},
			{[]Play{
				{"P2", S7},
				{"P3", S10},
				{"P4", SK},
				{"P1", HQ},
			}, "P1"},
			{[]Play{
				{"P1", S8},
				{"P2", HK},
				{"P3", SJ},
				{"P4", SA},
			}, "P2"},
		},
	}
}

func newGameWithBelote() Game {
	game := newNormalGame()
	game.Turns = game.Turns[:len(game.Turns)-1]
	game.Turns = append(game.Turns, Turn{[]Play{
		{"P1", HK}, // belote with H_Q
		{"P2", S8},
		{"P3", SJ},
		{"P4", SA},
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
	game.Turns = []Turn{
		{[]Play{
			{"P1", C7},
			{"P2", C10},
			{"P3", CK},
			{"P4", H9},
		}, "P2"},
		{[]Play{
			{"P4", D8},
			{"P1", DJ},
			{"P2", DK},
			{"P3", D7},
		}, "P2"},
		{[]Play{
			{"P2", CJ},
			{"P3", CA},
			{"P4", H10},
			{"P1", C8},
		}, "P3"},
		{[]Play{
			{"P4", D9},
			{"P1", DQ},
			{"P2", DA},
			{"P3", H7},
		}, "P2"},
		{[]Play{
			{"P3", H8},
			{"P4", D10},
			{"P1", HJ},
			{"P2", HA},
		}, "P2"},
		{[]Play{
			{"P1", C9},
			{"P2", CQ},
			{"P3", S9},
			{"P4", SQ},
		}, "P2"},
		{[]Play{
			{"P2", S7},
			{"P3", S10},
			{"P4", SK},
			{"P1", HQ},
		}, "P3"},
		{[]Play{
			{"P1", S8},
			{"P2", HK},
			{"P3", SJ},
			{"P4", SA},
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
	game.Turns = []Turn{
		{[]Play{
			{"P1", C7},
			{"P2", C10},
			{"P3", CK}, // first belote with C_Q
			{"P4", H9},
		}, "P2"},
		{[]Play{
			{"P4", D8},
			{"P1", DJ},
			{"P2", DK}, // could be a belote but too late
			{"P3", D7},
		}, "P1"},
		{[]Play{
			{"P2", CJ},
			{"P3", CA},
			{"P4", H10},
			{"P1", C8},
		}, "P2"},
		{[]Play{
			{"P4", D9},
			{"P1", DA},
			{"P2", DQ}, // could be a rebelote but too late
			{"P3", H7},
		}, "P4"},
		{[]Play{
			{"P3", H8},
			{"P4", D10},
			{"P1", HJ},
			{"P2", HA},
		}, "P2"},
		{[]Play{
			{"P1", C9},
			{"P2", CQ}, // rebelote
			{"P3", S9},
			{"P4", SQ},
		}, "P1"},
		{[]Play{
			{"P2", S7},
			{"P3", S10},
			{"P4", SK},
			{"P1", HQ},
		}, "P3"},
		{[]Play{
			{"P1", S8},
			{"P2", HK},
			{"P3", SJ},
			{"P4", SA},
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
	game.Turns = []Turn{
		{[]Play{
			{"P1", C7},
			{"P2", C10},
			{"P3", CK},
			{"P4", H9},
		}, "P4"},
		{[]Play{
			{"P4", D8},
			{"P1", DJ},
			{"P2", DK},
			{"P3", D7},
		}, "P2"},
		{[]Play{
			{"P2", CJ},
			{"P3", CA},
			{"P4", H10},
			{"P1", C8},
		}, "P4"},
		{[]Play{
			{"P4", D9},
			{"P1", DQ},
			{"P2", H7},
			{"P3", DA},
		}, "P2"},
		{[]Play{
			{"P3", H8},
			{"P4", D10},
			{"P1", HA},
			{"P2", HJ},
		}, "P2"},
		{[]Play{
			{"P1", C9},
			{"P2", CQ},
			{"P3", S9},
			{"P4", SQ},
		}, "P2"},
		{[]Play{
			{"P1", S8},
			{"P2", HK},
			{"P3", SJ},
			{"P4", SA},
		}, "P2"},
		{[]Play{
			{"P2", S7},
			{"P3", S10},
			{"P4", SK},
			{"P1", HQ},
		}, "P1"},
	}
	return game
}

func newGameWithCapotWon() Game {
	game := newGameWithCapotLost()
	game.Turns = game.Turns[:len(game.Turns)-1]
	game.Turns = append(game.Turns, Turn{[]Play{
		{"P2", S7},
		{"P3", S10},
		{"P4", HQ},
		{"P1", SK},
	}, "P4"},
	)
	return game
}

func TestCountingPhase(test *testing.T) {
	assert := assert.New(test)

	test.Run("should go to counting phase on last game", func(test *testing.T) {
		game := newNormalGame()
		game.Phase = Playing
		game.Turns = game.Turns[:len(game.Turns)-1]
		game.Turns = append(game.Turns, Turn{[]Play{
			{"P1", S8},
			{"P2", HK},
			{"P3", SJ},
		}, ""},
		)

		err := game.Play("P4", SA)
		assert.NoError(err)

		assert.Equal(Counting, game.Phase)
	})
}

func TestPlayersCards(test *testing.T) {
	assert := assert.New(test)

	test.Run("should count correctly in a game", func(test *testing.T) {
		game := newNormalGame()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID{H8, D10, HJ, HA, S7, S10, SK, HQ}, playerCards["P1"])
		assert.Equal([]CardID{D8, DJ, DK, D7, C9, CQ, S9, SQ, S8, HK, SJ, SA}, playerCards["P2"])
		assert.Equal([]CardID{D9, DQ, DA, H7}, playerCards["P3"])
		assert.Equal([]CardID{C7, C10, CK, H9, CJ, CA, H10, C8}, playerCards["P4"])
	})

	test.Run("should count correctly in a game with BELOTE", func(test *testing.T) {
		game := newGameWithBelote()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID{H8, D10, HJ, HA, S7, S10, SK, HQ, HK, S8, SJ, SA}, playerCards["P1"])
		assert.Equal([]CardID{D8, DJ, DK, D7, C9, CQ, S9, SQ}, playerCards["P2"])
		assert.Equal([]CardID{D9, DQ, DA, H7}, playerCards["P3"])
		assert.Equal([]CardID{C7, C10, CK, H9, CJ, CA, H10, C8}, playerCards["P4"])
	})

	test.Run("should count correctly in a game with NO-TRUMP", func(test *testing.T) {
		game := newGameWithNoTrump()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID(nil), playerCards["P1"])
		assert.Equal([]CardID{C7, C10, CK, H9, D8, DJ, DK, D7, D9, DQ, DA, H7, H8, D10, HJ, HA, C9, CQ, S9, SQ}, playerCards["P2"])
		assert.Equal([]CardID{CJ, CA, H10, C8, S7, S10, SK, HQ}, playerCards["P3"])
		assert.Equal([]CardID{S8, HK, SJ, SA}, playerCards["P4"])
	})

	test.Run("should count correctly in a game with ALL-TRUMP and BELOTE (for odd team)", func(test *testing.T) {
		game := newGameWithAllTrump()

		playerCards := game.getPlayersCards()
		assert.Equal([]CardID{D8, DJ, DK, D7, C9, CQ, S9, SQ}, playerCards["P1"])
		assert.Equal([]CardID{C7, C10, CK, H9, CJ, CA, H10, C8, H8, D10, HJ, HA}, playerCards["P2"])
		assert.Equal([]CardID{S7, S10, SK, HQ, S8, HK, SJ, SA}, playerCards["P3"])
		assert.Equal([]CardID{D9, DA, DQ, H7}, playerCards["P4"])
	})
}

func TestCounting(test *testing.T) {
	assert := assert.New(test)

	test.Run("should count correctly in a LOST game", func(test *testing.T) {
		game := newNormalGame()

		game.calculatesTeamPointsAndScores()
		assert.Equal(72, game.Points["odd"])
		assert.Equal(90, game.Points["even"])
		assert.Equal(72, game.Scores["odd"])
		assert.Equal(160+80, game.Scores["even"])

		assert.Equal(162, game.Points["odd"]+game.Points["even"])
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
		assert.Equal(99+20, game.Points["odd"])
		assert.Equal(63, game.Points["even"])
		assert.Equal((80+160+20)*2, game.Scores["odd"])
		assert.Equal(0, game.Scores["even"])

		assert.Equal(182, game.Points["odd"]+game.Points["even"])
	})

	test.Run("should count correctly in a WON game with BELOTE", func(test *testing.T) {
		game := newGameWithBelote()

		game.calculatesTeamPointsAndScores()
		assert.Equal(99+20, game.Points["odd"])
		assert.Equal(63, game.Points["even"])
		assert.Equal((80 + 99 + 20), game.Scores["odd"])
		assert.Equal(63, game.Scores["even"])

		assert.Equal(182, game.Points["odd"]+game.Points["even"])
	})

	test.Run("should count correctly in a LOST game with NO-TRUMP", func(test *testing.T) {
		game := newGameWithNoTrump()

		game.calculatesTeamPointsAndScores()
		assert.Equal(49, game.Points["odd"])
		assert.Equal(113, game.Points["even"])
		assert.Equal(49, game.Scores["odd"])
		assert.Equal(160+80, game.Scores["even"])

		assert.Equal(162, game.Points["odd"]+game.Points["even"])
	})

	test.Run("should count correctly in a WON game with ALL-TRUMP and BELOTE (for odd team)", func(test *testing.T) {
		game := newGameWithAllTrump()

		game.calculatesTeamPointsAndScores()
		assert.Equal(76, game.Points["odd"])
		assert.Equal(86, game.Points["even"])
		assert.Equal(20+76, game.Scores["odd"]) // the belote was taken by the odd team but the even team had taken the bid
		assert.Equal(80+86, game.Scores["even"])

		assert.Equal(162, game.Points["odd"]+game.Points["even"])
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
		assert.Equal(76, game.Points["odd"])
		assert.Equal(86, game.Points["even"])
		assert.Equal(20*4, game.Scores["odd"]) // the belote was taken by the odd team but the even team had taken the bid
		assert.Equal((160+80)*4, game.Scores["even"])
	})

	test.Run("should count correctly with CAPOT LOST", func(test *testing.T) {
		game := newGameWithCapotLost()

		game.calculatesTeamPointsAndScores()
		assert.Equal(27, game.Points["odd"])
		assert.Equal(135, game.Points["even"])
		assert.Equal(320, game.Scores["odd"])
		assert.Equal(0, game.Scores["even"])
	})

	test.Run("should count correctly with CAPOT WON", func(test *testing.T) {
		game := newGameWithCapotWon()

		game.calculatesTeamPointsAndScores()
		assert.Equal(0, game.Points["odd"])
		assert.Equal(162, game.Points["even"])
		assert.Equal(0, game.Scores["odd"])
		assert.Equal(500, game.Scores["even"])
	})
}

func TestRestarting(test *testing.T) {
	assert := assert.New(test)

	test.Run("should be able to restart by reinitializing everything but the score", func(test *testing.T) {
		game := newNormalGame()
		game.calculatesTeamPointsAndScores()

		err := game.Start()
		if err != nil {
			test.Fatal(err)
		}

		assert.Equal(Bidding, game.Phase)
		assert.Equal(0, len(game.Bids))
		assert.Equal(0, len(game.Turns))
		assert.Equal(0, len(game.Points))
		assert.Equal(0, len(game.Deck))

		assert.Equal(72, game.Scores["odd"])
		assert.Equal(160+80, game.Scores["even"])

		assert.Equal(8, len(game.Players["P1"].Hand))
		assert.Equal(8, len(game.Players["P2"].Hand))
		assert.Equal(8, len(game.Players["P3"].Hand))
		assert.Equal(8, len(game.Players["P4"].Hand))

		assert.Equal(1, game.Players["P2"].Order)
		assert.Equal(1, game.Players["P2"].InitialOrder)
	})

	test.Run("should be able to add to score without reinitializing", func(test *testing.T) {
		game := newNormalGame()
		game.Scores["even"] = 1000
		game.Scores["odd"] = 2000

		game.calculatesTeamPointsAndScores()

		assert.Equal(72, game.Points["odd"])
		assert.Equal(90, game.Points["even"])

		assert.Equal(2000+72, game.Scores["odd"])
		assert.Equal(1000+160+80, game.Scores["even"])
	})
}
