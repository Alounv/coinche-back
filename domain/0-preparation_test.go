package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newGameWith2Players() Game {
	return Game{
		ID:      2,
		Name:    "GAME TWO",
		Players: map[string]Player{"P1": {}, "P2": {}},
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
