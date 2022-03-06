package domain

import (
	"coinche/utilities"
	"errors"
	"fmt"
)

const (
	ErrNotPlaying    = "NOT IN PLAYING PHASE"
	ErrCardNotInHand = "CARD NOT IN HAND"
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
	game.deck = []int{}
}

func (game *Game) draw(order int) []int {
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
	hand := []int{}

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

func (game *Game) Play(playerName string, card int) error {
	if game.Phase != Playing {
		return errors.New(ErrNotPlaying)
	}

	var player Player
	for name, p := range game.Players {
		if name == playerName {
			player = p
		}
	}

	fmt.Println(playerName, player)
	err := game.checkPlayerTurn(playerName)
	if err != nil {
		return err
	}

	for _, index := range player.Hand {
		if index == card {
			player.Hand = utilities.Remove(player.Hand, card)
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
