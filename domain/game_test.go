package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type cardTest struct {
	card      cardID
	trump     Color
	firstCard cardID
	strength  Strength
}

func TestGetCardValue(test *testing.T) {
	assert := assert.New(test)

	cardTests := []cardTest{
		{S_10, NoTrump, S_10, Ten},
		{S_10, NoTrump, S_9, Ten},
		{S_10, NoTrump, C_9, 0},
		{S_10, AllTrump, C_9, TTen},
		{S_10, Spade, C_9, TTen},
		{C_8, NoTrump, C_8, Eight},
		{C_8, NoTrump, C_9, Eight},
		{C_8, NoTrump, H_9, 0},
		{C_8, AllTrump, H_9, TEight},
		{C_8, Club, H_9, TEight},
	}

	for _, t := range cardTests {
		test.Run("value should be correct", func(test *testing.T) {
			got := getCardValue(t.card, t.trump, t.firstCard)

			assert.Equal(t.strength, got)
		})
	}
}
