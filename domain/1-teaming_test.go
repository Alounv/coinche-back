package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		assert.NoError(testGame.CanStartBidding())
	})

	test.Run("should not be ready with on team of one", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {Team: "A"}, "P2": {Team: "A"}, "P3": {}, "P4": {Team: "B"}},
			Phase:   Teaming,
		}
		assert.Error(testGame.CanStartBidding())
	})

	test.Run("should not be ready with on team of three", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {Team: "A"}, "P2": {Team: "A"}, "P3": {Team: "A"}, "P4": {Team: "B"}},
			Phase:   Teaming,
		}
		assert.Error(testGame.CanStartBidding())
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

		err := testGame.StartBidding()

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
		err := testGame.StartBidding()

		assert.NoError(err)
		assert.Equal(Bidding, testGame.Phase)
		assert.Equal(4, testGame.Players["P1"].Order)
		assert.Equal(1, testGame.Players["P2"].Order)
		assert.Equal(2, testGame.Players["P3"].Order)
		assert.Equal(3, testGame.Players["P4"].Order)
	})
}
