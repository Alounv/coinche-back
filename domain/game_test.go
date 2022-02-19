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
			Players: map[string]Player{"P1": {"A"}, "P2": {"A"}, "P3": {"B"}, "P4": {"B"}},
			Phase:   Teaming,
		}
		assert.NoError(testGame.CanStart())
	})

	test.Run("should not be ready with on team of one", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {"A"}, "P2": {"A"}, "P3": {}, "P4": {"B"}},
			Phase:   Teaming,
		}
		assert.Error(testGame.CanStart())
	})

	test.Run("should not be ready with on team of three", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {"A"}, "P2": {"A"}, "P3": {"A"}, "P4": {"B"}},
			Phase:   Teaming,
		}
		assert.Error(testGame.CanStart())
	})
}

func TestStart(test *testing.T) {
	assert := assert.New(test)

	test.Run("should start when the game can start", func(test *testing.T) {
		testGame := Game{
			ID:      2,
			Name:    "GAME TWO",
			Players: map[string]Player{"P1": {"A"}, "P2": {"A"}, "P3": {"B"}, "P4": {"B"}},
			Phase:   Teaming,
		}

		err := testGame.Start()

		assert.NoError(err)
		assert.Equal(Bidding, testGame.Phase)
	})
}
