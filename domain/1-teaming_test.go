package domain

import (
	"coinche/utilities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newGameWith2Players() Game {
	return Game{
		ID:      2,
		Name:    "GAME TWO",
		Players: map[string]Player{"P1": {}, "P2": {}},
		Phase:   Teaming,
	}
}

func newGameWith4Players() Game {
	return Game{
		ID:      2,
		Name:    "GAME TWO",
		Players: map[string]Player{"P1": {}, "P2": {}, "P3": {}, "P4": {}},
		Phase:   Teaming,
	}
}

func TestAddPlayer(test *testing.T) {
	assert := assert.New(test)

	test.Run("full game should be full", func(test *testing.T) {
		game := newGameWith4Players()

		got := game.IsFull()

		assert.Equal(true, got)
	})

	test.Run("game not full should not be full", func(test *testing.T) {
		game := newGameWith2Players()

		got := game.IsFull()

		assert.Equal(false, got)
	})

	test.Run("should add player", func(test *testing.T) {
		game := newGameWith2Players()

		err := game.AddPlayer("P3")

		assert.Equal(3, len(game.Players))
		assert.NoError(err)
	})

	test.Run("should fail to add player when full", func(test *testing.T) {
		game := newGameWith4Players()

		err := game.AddPlayer("P5")

		assert.Equal(err.Error(), ErrGameFull)
	})

	test.Run("should fail to add player already in game", func(test *testing.T) {
		game := newGameWith2Players()

		err := game.AddPlayer("P2")

		assert.Equal(2, len(game.Players))
		assert.Equal(err.Error(), ErrAlreadyInGame)
	})

	test.Run("should return AlreadyInGame error if game is full and player is already in game", func(test *testing.T) {
		game := newGameWith4Players()

		err := game.AddPlayer("P3")

		assert.Equal(err.Error(), ErrAlreadyInGame)
	})
}

func TestRemovePlayer(test *testing.T) {
	assert := assert.New(test)

	test.Run("should remove player", func(test *testing.T) {
		game := newGameWith4Players()

		err := game.RemovePlayer("P2")

		assert.NoError(err)
		assert.Equal(3, len(game.Players))
	})

	test.Run("should fail to remove player not in game", func(test *testing.T) {
		game := newGameWith2Players()

		err := game.RemovePlayer("P3")

		assert.Error(err)
	})
}

func TestTeamingPhase(test *testing.T) {
	assert := assert.New(test)

	test.Run("should be in teaming phase", func(test *testing.T) {
		game := newGameWith4Players()

		assert.Equal(Teaming, game.Phase)
	})

	test.Run("can create a team", func(test *testing.T) {
		game := newGameWith4Players()

		err := game.AssignTeam("P1", "Team1")

		assert.NoError(err)
		assert.Equal("Team1", game.Players["P1"].Team)
	})

	test.Run("can join a team", func(test *testing.T) {
		game := newGameWith4Players()
		game.Players = map[string]Player{
			"P1": {Team: "Team1"},
			"P2": {},
			"P3": {},
			"P4": {},
		}

		err := game.AssignTeam("P2", "Team1")

		assert.NoError(err)
		assert.Equal("Team1", game.Players["P1"].Team)
		assert.Equal("Team1", game.Players["P2"].Team)
	})

	test.Run("should fail when team is full", func(test *testing.T) {
		game := newGameWith4Players()
		game.Players = map[string]Player{
			"P1": {Team: "Team1"},
			"P2": {Team: "Team1"},
			"P3": {},
			"P4": {},
		}

		err := game.AssignTeam("P3", "Team1")

		assert.Equal(err.Error(), ErrTeamFull)
	})

	test.Run("can leave a team", func(test *testing.T) {
		game := newGameWith4Players()
		game.Players = map[string]Player{
			"P1": {Team: "Team1"},
			"P2": {Team: "Team1"},
			"P3": {},
			"P4": {},
		}

		err := game.ClearTeam("P2")

		assert.NoError(err)
		assert.Equal("", game.Players["P2"].Team)
	})
}

func newTeamingGameWithoutPlayers() Game {
	return Game{
		ID:    2,
		Name:  "GAME TWO",
		Phase: Teaming,
		Deck:  NewDeck(),
	}
}

func TestCanStart(test *testing.T) {
	assert := assert.New(test)

	test.Run("should be ready to start with two teams of two", func(test *testing.T) {
		game := newTeamingGameWithoutPlayers()
		game.Players = map[string]Player{
			"P1": {Team: "A"},
			"P2": {Team: "A"},
			"P3": {Team: "B"},
			"P4": {Team: "B"},
		}

		assert.NoError(game.CanStartBidding())
	})

	test.Run("should not be ready with on team of one", func(test *testing.T) {
		game := newTeamingGameWithoutPlayers()
		game.Players = map[string]Player{
			"P1": {Team: "A"},
			"P2": {Team: "A"},
			"P3": {},
			"P4": {Team: "B"},
		}

		assert.Error(game.CanStartBidding())
	})

	test.Run("should not be ready with on team of three", func(test *testing.T) {
		game := newTeamingGameWithoutPlayers()
		game.Players = map[string]Player{
			"P1": {Team: "A"},
			"P2": {Team: "A"},
			"P3": {Team: "A"},
			"P4": {Team: "B"},
		}

		assert.Error(game.CanStartBidding())
	})
}

func TestStart(test *testing.T) {
	assert := assert.New(test)

	test.Run("should start when the game can start", func(test *testing.T) {
		game := newTeamingGameWithoutPlayers()
		game.Players = map[string]Player{
			"P1": {Team: "A"},
			"P2": {Team: "A"},
			"P3": {Team: "B"},
			"P4": {Team: "B"},
		}

		err := game.StartBidding()

		assert.NoError(err)
		assert.Equal(Bidding, game.Phase)
		assert.Equal(2, utilities.Abs(game.Players["P1"].Order-game.Players["P2"].Order))
		assert.Equal(2, utilities.Abs(game.Players["P3"].Order-game.Players["P4"].Order))
	})

	test.Run("should rotate order when order exists", func(test *testing.T) {
		game := newTeamingGameWithoutPlayers()
		game.Players = map[string]Player{
			"P1": {Team: "odd", InitialOrder: 1},
			"P2": {Team: "even", InitialOrder: 2},
			"P3": {Team: "odd", InitialOrder: 3},
			"P4": {Team: "even", InitialOrder: 4},
		}

		game.Phase = Teaming
		err := game.StartBidding()

		assert.NoError(err)
		assert.Equal(Bidding, game.Phase)
		assert.Equal(4, game.Players["P1"].Order)
		assert.Equal(1, game.Players["P2"].Order)
		assert.Equal(2, game.Players["P3"].Order)
		assert.Equal(3, game.Players["P4"].Order)
	})
}
