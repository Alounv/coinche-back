package domain

import (
	"errors"
)

const (
	ErrNotPlaying            = "NOT IN PLAYING PHASE"
	ErrCardNotInHand         = "CARD NOT IN HAND"
	ErrShouldPlayAskedColor  = "SHOULD PLAY ASKED COLOR"
	ErrShouldPlayBiggerTrump = "SHOULD PLAY BIGGER TRUMP"
	ErrShouldPlayTrump       = "SHOULD PLAY TRUMP"
)

func (game *Game) startPlaying() {
	game.Phase = Playing
	game.setTrump()
	game.resetOrderAsInitialOrder()
	game.distributeCards()
}

func (game *Game) setTrump() {
	var maxValue BidValue
	for value := range game.Bids {
		if value > maxValue {
			maxValue = value
		}
	}

	lastBid := game.Bids[maxValue]

	game.trump = lastBid.Color
}

func (game *Game) distributeCards() {
	for name, player := range game.Players {
		player.Hand = game.draw(player.Order)
		game.Players[name] = player
	}
	game.deck = []cardID{}
}

func (game *Game) draw(order int) []cardID {
	base := order - 1

	deckIndexes := []int{
		base*3 + 0,
		base*3 + 1,
		base*3 + 2,

		12 + base*2 + 0,
		12 + base*2 + 1,

		20 + base*3 + 0,
		20 + base*3 + 1,
		20 + base*3 + 2,
	}
	hand := []cardID{}

	for _, deckIndex := range deckIndexes {
		hand = append(hand, game.deck[deckIndex])
	}

	return hand
}

func (game *Game) isNewTurn() bool {
	if len(game.turns) == 0 {
		return true
	}
	lastTurn := game.turns[len(game.turns)-1]
	return len(lastTurn.plays) >= 4
}

func (turn turn) askedColor() Color {
	firstPlay := turn.plays[0]
	return cards[firstPlay.card].color
}

func hasColor(color Color, hand []cardID) bool {
	for _, card := range hand {
		if cards[card].color == color {
			return true
		}
	}
	return false
}

func hasNoTrump(trump Color, hand []cardID) bool {
	for _, card := range hand {
		if cards[card].color == trump {
			return false
		}
	}
	return true
}

func (turn turn) getBiggestTrumpStrength(trump Color) Strength {
	var biggestTrumpStrength Strength
	for _, play := range turn.plays {
		card := cards[play.card]
		color := card.color
		strength := card.TrumpStrength
		if color == trump {
			if strength > biggestTrumpStrength {
				biggestTrumpStrength = strength
			}
		}
	}
	return biggestTrumpStrength
}

func (turn turn) isTheBiggestTrump(card cardID, trump Color) bool {
	isTrump := cards[card].color == trump

	if isTrump && cards[card].TrumpStrength > turn.getBiggestTrumpStrength(trump) {
		return true
	}

	return false
}

func (turn turn) hasNoBiggerTrump(hand []cardID, trump Color) bool {
	for _, cardID := range hand {
		card := cards[cardID]
		if card.color == trump && card.TrumpStrength > turn.getBiggestTrumpStrength(trump) {
			return false
		}
	}
	return true
}

func (game *Game) canPlayCard(card cardID, playerName string) error {
	player := game.Players[playerName]

	if len(game.turns) == 0 {
		return nil
	}

	lastTurn := game.turns[len(game.turns)-1]
	playCount := len(lastTurn.plays)

	if playCount == 0 {
		return nil
	}

	askedColor := lastTurn.askedColor()
	color := cards[card].color

	if color == askedColor {
		return nil
	}

	if hasColor(askedColor, player.Hand) {
		return errors.New(ErrShouldPlayAskedColor)
	}

	if hasNoTrump(game.trump, player.Hand) {
		return nil
	}

	if playCount > 1 {
		winnerTeam := game.Players[lastTurn.getWinner(game.trump)].Team
		if winnerTeam == player.Team {
			return nil
		}
	}

	if color != game.trump {
		return errors.New(ErrShouldPlayTrump)
	}

	if lastTurn.isTheBiggestTrump(card, game.trump) || lastTurn.hasNoBiggerTrump(player.Hand, game.trump) {
		return nil
	}

	return errors.New(ErrShouldPlayBiggerTrump)
}

func (game *Game) Play(playerName string, card cardID) error {
	if game.Phase != Playing {
		return errors.New(ErrNotPlaying)
	}

	var player Player
	for name, p := range game.Players {
		if name == playerName {
			player = p
		}
	}

	err := game.checkPlayerTurn(playerName)
	if err != nil {
		return err
	}

	err = game.canPlayCard(card, playerName)
	if err != nil {
		return err
	}

	for _, cardID := range player.Hand {
		if cardID == card {
			player.Hand = removeCard(player.Hand, card)
			game.Players[playerName] = player

			newPlay := play{
				playerName: playerName,
				card:       card,
			}

			game.rotateOrder()

			if game.isNewTurn() {
				game.turns = append(game.turns, turn{
					plays: []play{newPlay},
				})
				return nil
			}

			lastTurnIndex := len(game.turns) - 1
			lastTurn := game.turns[lastTurnIndex]

			lastTurn.plays = append(lastTurn.plays, newPlay)

			if len(lastTurn.plays) == 4 {
				lastTurn.setWinner(Color(game.trump))
				game.setFirstPlayer(lastTurn.winner)
			}

			game.turns[lastTurnIndex] = lastTurn

			return nil
		}
	}

	return errors.New(ErrCardNotInHand)
}

func removeCard(slice []cardID, valueToRemove cardID) []cardID {
	for index, value := range slice {
		if value == valueToRemove {
			return removeWithIndex(slice, index)
		}
	}
	return slice
}

func removeWithIndex(slice []cardID, index int) []cardID {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
