package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPlayer(test *testing.T) {
	assert := assert.New(test)
	testGame := Game{
		ID:      1,
		Name:    "GameName",
		Players: map[string]Player{"P1": {}, "P2": {}},
	}
	fullTestGame := Game{
		ID:      2,
		Name:    "FullGameName",
		Players: map[string]Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}},
	}

	test.Run("full game should be full", func(test *testing.T) {
		got := fullTestGame.IsFull()
		assert.Equal(true, got)
	})

	test.Run("game not full should not be full", func(test *testing.T) {
		got := testGame.IsFull()

		assert.Equal(false, got)
	})

	test.Run("should add player", func(test *testing.T) {
		err := testGame.AddPlayer("P3")

		assert.Equal(testGame.Players, map[string]Player{"P1": {}, "P2": {}, "P3": {}})
		assert.NoError(err)
	})

	test.Run("should fail to add player when full", func(test *testing.T) {
		err := fullTestGame.AddPlayer("P5")

		assert.Equal(err.Error(), ErrGameFull)
	})

	test.Run("should fail to add player already in game", func(test *testing.T) {
		err := testGame.AddPlayer("P2")

		assert.Equal(map[string]Player{"P1": {}, "P2": {}, "P3": {}}, testGame.Players)
		assert.Equal(err.Error(), ErrAlreadyInGame)
	})

	test.Run("should return AlreadyInGame error if game is full and player is already in game", func(test *testing.T) {
		err := fullTestGame.AddPlayer("P3")

		assert.Equal(err.Error(), ErrAlreadyInGame)
	})
}

func TestRemovePlayer(test *testing.T) {
	assert := assert.New(test)
	testGame := Game{
		ID:      2,
		Name:    "FullGameName",
		Players: map[string]Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}},
	}

	test.Run("should remove player", func(test *testing.T) {
		err := testGame.RemovePlayer("P2")

		assert.NoError(err)
		assert.Equal(map[string]Player{"P1": {}, "P3": {}, "P4": {}}, testGame.Players)
	})

	test.Run("should fail to remove player not in game", func(test *testing.T) {
		err := testGame.RemovePlayer("P2")

		assert.Error(err)
	})
}

func TestGamePhases(test *testing.T) {
	assert := assert.New(test)
	testGame := Game{
		ID:      2,
		Name:    "GAME TWO",
		Players: map[string]Player{"P1": {}, "P2": {}, "P3": {}},
	}

	test.Run("should be in preparation phase", func(test *testing.T) {
		assert.Equal(Preparation, testGame.Phase)
	})

	test.Run("should be in teaming phase when full", func(test *testing.T) {
		err := testGame.AddPlayer("P4")
		assert.NoError(err)
		assert.Equal(Teaming, testGame.Phase)
	})

	test.Run("should stay in teaming if trying to add existing player", func(test *testing.T) {
		err := testGame.AddPlayer("P4")
		assert.Error(err)
		assert.Equal(Teaming, testGame.Phase)
	})

	test.Run("should go in pause phase is players go less than 4 after preparation", func(test *testing.T) {
		err := testGame.RemovePlayer("P4")
		assert.NoError(err)
		assert.Equal(Pause, testGame.Phase)
	})
}

func TestTeamingPhase(test *testing.T) {
	assert := assert.New(test)
	testGame := Game{
		ID:      2,
		Name:    "GAME TWO",
		Players: map[string]Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}},
		Phase:   Teaming,
	}

	test.Run("should be in teaming phase", func(test *testing.T) {
		assert.Equal(Teaming, testGame.Phase)
	})

	test.Run("can create a team", func(test *testing.T) {
		err := testGame.AssignTeam("P1", "Team1")

		assert.NoError(err)
		assert.Equal("Team1", testGame.Players["P1"].Team)
	})

	test.Run("can join a team", func(test *testing.T) {
		err := testGame.AssignTeam("P2", "Team1")

		assert.NoError(err)
		assert.Equal("Team1", testGame.Players["P1"].Team)
		assert.Equal("Team1", testGame.Players["P2"].Team)
	})

	test.Run("should fail when team is full", func(test *testing.T) {
		err := testGame.AssignTeam("P3", "Team1")

		assert.Equal(err.Error(), ErrTeamFull)
	})

	test.Run("can leave a team", func(test *testing.T) {
		err := testGame.ClearTeam("P2")

		assert.NoError(err)
		assert.Equal("", testGame.Players["P2"].Team)
	})
}

func TestCanStart(test *testing.T) {
	assert := assert.New(test)

	test.Run("should be ready to start with two teams of two", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {Team: "A"}, "P2": {Team: "A"}, "P3": {Team: "B"}, "P4": {Team: "B"}},
			Phase:   Teaming,
		}
		assert.NoError(testGame.CanStart())
	})

	test.Run("should not be ready with on team of one", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {Team: "A"}, "P2": {Team: "A"}, "P3": {}, "P4": {Team: "B"}},
			Phase:   Teaming,
		}
		assert.Error(testGame.CanStart())
	})

	test.Run("should not be ready with on team of three", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {Team: "A"}, "P2": {Team: "A"}, "P3": {Team: "A"}, "P4": {Team: "B"}},
			Phase:   Teaming,
		}
		assert.Error(testGame.CanStart())
	})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestStart(test *testing.T) {
	assert := assert.New(test)

	test.Run("should start when the game can start", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {Team: "A"}, "P2": {Team: "A"}, "P3": {Team: "B"}, "P4": {Team: "B"}},
			Phase:   Teaming,
		}

		err := testGame.Start()

		assert.NoError(err)
		assert.Equal(Bidding, testGame.Phase)
		assert.Equal(2, abs(testGame.Players["P1"].Order-testGame.Players["P2"].Order))
		assert.Equal(2, abs(testGame.Players["P3"].Order-testGame.Players["P4"].Order))
	})

	test.Run("should rotate order when order exists", func(test *testing.T) {
		testGame := Game{
			ID:   2,
			Name: "GAME TWO",
			Players: map[string]Player{
				"P1": {Team: "odd", InitialOrder: 1},
				"P2": {Team: "even", InitialOrder: 2},
				"P3": {Team: "odd", InitialOrder: 3},
				"P4": {Team: "even", InitialOrder: 4},
			},
			Phase: Teaming,
		}

		testGame.Phase = Teaming
		err := testGame.Start()

		assert.NoError(err)
		assert.Equal(Bidding, testGame.Phase)
		assert.Equal(4, testGame.Players["P1"].Order)
		assert.Equal(1, testGame.Players["P2"].Order)
		assert.Equal(2, testGame.Players["P3"].Order)
		assert.Equal(3, testGame.Players["P4"].Order)
	})
}

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
		Bids:  map[Value]Bid{},
	}
	test.Run("should fail if not in bidding", func(test *testing.T) {
		teamingGame := Game{
			ID:    2,
			Phase: Teaming,
		}
		err := teamingGame.PlaceBid("P1", Eight, Spade)

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
		err := testGame.PlaceBid("P1", Eight, Spade)

		assert.NoError(err)
		assert.Equal(want, testGame.Bids[Eight])
	})

	test.Run("should place another bid", func(test *testing.T) {
		want := Bid{Player: "P2", Color: Club}
		err := testGame.PlaceBid("P2", Nine, Club)

		assert.NoError(err)
		assert.Equal(want, testGame.Bids[Nine])
	})

	test.Run("should fail if placing a bid smaller or equal to previous bid", func(test *testing.T) {
		err := testGame.PlaceBid("P3", Eight, Club)

		assert.Error(err)
		assert.Equal(ErrBidTooSmall, err.Error())
	})

	test.Run("should fail if the player is not the right one", func(test *testing.T) {
		err := testGame.PlaceBid("P1", Ten, Club)

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
		err := testGame.PlaceBid("P2", Eleven, Club)

		assert.Error(err)
		assert.Equal(ErrBiddingItsOwnColor, err.Error())
	})

	test.Run("same team player should not be able to coinche", func(test *testing.T) {
		err := testGame.PlaceBid("P2", Eleven, Heart)
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
		err := testGame.Coinche("P1")
		assert.NoError(err)

		err = testGame.Coinche("P4")
		assert.NoError(err)
	})

	test.Run("should start playing after two passes after coinche", func(t *testing.T) {
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
			Bids: map[Value]Bid{
				Eight: {Player: "P4", Color: Spade, Coinche: 2},
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
			Bids: map[Value]Bid{
				Eight: {Player: "P4", Color: Spade},
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
