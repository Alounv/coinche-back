package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type cardTest struct {
	card      CardID
	trump     Color
	firstCard CardID
	strength  Strength
}

func TestGetCardValue(test *testing.T) {
	assert := assert.New(test)

	cardTests := []cardTest{
		{S10, NoTrump, S10, Ten},
		{S10, NoTrump, S9, Ten},
		{S10, NoTrump, C9, 0},
		{S10, AllTrump, C9, 0},
		{S10, Spade, C9, TTen},
		{C8, NoTrump, C8, Eight},
		{C8, NoTrump, C9, Eight},
		{C8, NoTrump, H9, 0},
		{C8, AllTrump, H9, 0},
		{C8, Club, H9, TEight},
		{C8, AllTrump, C9, TEight},
	}

	for _, t := range cardTests {
		test.Run("value should be correct", func(test *testing.T) {
			got := getCardValue(t.card, t.trump, t.firstCard)

			assert.Equal(t.strength, got)
		})
	}
}
