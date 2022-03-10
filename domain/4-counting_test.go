package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
C_7, C_8, C_9, D_J, D_Q, H_J, H_Q, H_K,
C_10, C_J, C_Q, D_K, D_A, H_A, S_7, S_8,
C_K, C_A, D_7, H_7, H_8, S_9, S_10, S_J,
D_8, D_9, D_10, H_9, H_10, S_Q, S_K, S_A,
*/

/*
, , ,
, , ,
, , S_J,
, , S_A,
*/
func newLastPlayingGame() Game {
	return Game{
		ID:   2,
		Name: "GAME TWO",
		Players: map[string]Player{
			"P1": {Team: "odd", Order: 2, InitialOrder: 1, Hand: []cardID{}},
			"P2": {Team: "even", Order: 3, InitialOrder: 2, Hand: []cardID{}},
			"P3": {Team: "odd", Order: 4, InitialOrder: 3, Hand: []cardID{}},
			"P4": {Team: "even", Order: 1, InitialOrder: 4, Hand: []cardID{S_A}},
		},
		Phase: Playing,
		Bids:  map[BidValue]Bid{},
		trump: Heart,
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
				{"P1", H_K},
				{"P2", S_8},
				{"P3", S_J},
			}, ""},
		},
	}
}

func TestCounting(test *testing.T) {
	assert := assert.New(test)

	test.Run("should count correctly in a normal game", func(test *testing.T) {
		game := newLastPlayingGame()

		err := game.Play("P4", S_A)
		assert.NoError(err)

		assert.Equal(Counting, game.Phase)

		playerCards := game.getPlayersCards()
		assert.Equal([]cardID{H_8, D_10, H_J, H_A, S_7, S_10, S_K, H_Q, H_K, S_8, S_J, S_A}, playerCards["P1"])
		assert.Equal([]cardID{D_8, D_J, D_K, D_7, C_9, C_Q, S_9, S_Q}, playerCards["P2"])
		assert.Equal([]cardID{D_9, D_Q, D_A, H_7}, playerCards["P3"])
		assert.Equal([]cardID{C_7, C_10, C_K, H_9, C_J, C_A, H_10, C_8}, playerCards["P4"])

		teamPoints := game.getTeamPoints()
		assert.Equal(119, teamPoints["odd"])
		assert.Equal(63, teamPoints["even"])

		assert.Equal(182, teamPoints["odd"]+teamPoints["even"])
	})
}

// COINCHE / SURCHOINCE
// NOTRUPM
// ALLTRUMP
// CAPOT
// BELOTE
// 10DEDER

// action to quit

// action to restart
