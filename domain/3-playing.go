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
	game.resetOrderAsInitialOrder()
}

func (game *Game) isNewTurn() bool {
	if len(game.Turns) == 0 {
		return true
	}
	lastTurn := game.Turns[len(game.Turns)-1]
	return len(lastTurn.Plays) >= 4
}

func (turn Turn) askedColor() Color {
	firstPlay := turn.Plays[0]
	return cards[firstPlay.Card].color
}

func (player Player) hasCard(card CardID) bool {
	for _, CardID := range player.Hand {
		if CardID == card {
			return true
		}
	}
	return false
}

func (player Player) hasColor(color Color) bool {
	for _, card := range player.Hand {
		if cards[card].color == color {
			return true
		}
	}
	return false
}

func (player Player) hasNoTrump(trump Color) bool {
	for _, card := range player.Hand {
		if cards[card].color == trump {
			return false
		}
	}
	return true
}

func (turn Turn) getBiggestTrumpStrength(trump Color) Strength {
	var biggestTrumpStrength Strength
	for _, play := range turn.Plays {
		card := cards[play.Card]
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

func (turn Turn) isTheBiggestTrump(card CardID, trump Color) bool {
	isTrump := cards[card].color == trump || trump == AllTrump

	if isTrump && cards[card].TrumpStrength > turn.getBiggestTrumpStrength(trump) {
		return true
	}

	return false
}

func (game Game) isPartnerWinner(team string) bool {
	lastTurn := game.Turns[len(game.Turns)-1]
	playCount := len(lastTurn.Plays)
	trump := game.trump()

	if playCount < 2 {
		return false
	}
	winnerTeam := game.Players[lastTurn.getWinner(trump)].Team
	return winnerTeam == team
}

func (player Player) hasBiggerTrump(trump Color, turn Turn) bool {
	for _, CardID := range player.Hand {
		card := cards[CardID]
		isTrump := card.color == trump || trump == AllTrump
		if isTrump && card.TrumpStrength > turn.getBiggestTrumpStrength(trump) {
			return true
		}
	}
	return false
}

func (game *Game) canPlayCard(card CardID, playerName string) error {
	player := game.Players[playerName]

	if len(game.Turns) == 0 {
		return nil
	}

	lastTurn := game.Turns[len(game.Turns)-1]
	playCount := len(lastTurn.Plays)

	if playCount == 0 || playCount == 4 {
		return nil
	}

	askedColor := lastTurn.askedColor()
	color := cards[card].color
	trump := game.trump()

	if player.hasColor(askedColor) {
		if color != askedColor {
			return errors.New(ErrShouldPlayAskedColor)
		}

		if askedColor != trump && trump != AllTrump {
			return nil
		}

		if !lastTurn.isTheBiggestTrump(card, trump) && player.hasBiggerTrump(trump, lastTurn) {
			return errors.New(ErrShouldPlayBiggerTrump)
		}

		return nil
	}

	// --- does not have asked color ---

	if player.hasNoTrump(trump) {
		return nil
	}

	if color == trump {
		return nil
	}

	isPartnerWinner := game.isPartnerWinner(player.Team)
	if isPartnerWinner {
		return nil
	}

	return errors.New(ErrShouldPlayTrump)
}

func (game *Game) createTurn(newPlay Play) {
	game.Turns = append(game.Turns, Turn{
		Plays: []Play{newPlay},
	})
}

func (game *Game) updateTurn(newPlay Play) {
	lastTurnIndex := len(game.Turns) - 1
	lastTurn := game.Turns[lastTurnIndex]

	lastTurn.Plays = append(lastTurn.Plays, newPlay)

	if len(lastTurn.Plays) == 4 {
		lastTurn.setWinner(game.trump())
		game.setFirstPlayer(lastTurn.Winner)
	}

	game.Turns[lastTurnIndex] = lastTurn
}

func (game *Game) allCardsPlayed() bool {
	return len(game.Turns) == 8 && len(game.Turns[7].Plays) == 4
}

func (game *Game) Play(playerName string, card CardID) error {
	if game.Phase != Playing {
		return errors.New(ErrNotPlaying)
	}

	err := game.checkPlayerTurn(playerName)
	if err != nil {
		return err
	}

	err = game.canPlayCard(card, playerName)
	if err != nil {
		return err
	}

	player := game.Players[playerName]

	if !player.hasCard(card) {
		return errors.New(ErrCardNotInHand)
	}

	player.Hand = removeCard(player.Hand, card)
	game.Players[playerName] = player

	game.rotateOrder()

	newPlay := Play{
		PlayerName: playerName,
		Card:       card,
	}

	if game.isNewTurn() {
		game.createTurn(newPlay)
		return nil
	}

	game.updateTurn(newPlay)

	if game.allCardsPlayed() {
		game.end()
	}

	return nil
}

func removeCard(slice []CardID, valueToRemove CardID) []CardID {
	for index, value := range slice {
		if value == valueToRemove {
			return removeWithIndex(slice, index)
		}
	}
	return slice
}

func removeWithIndex(slice []CardID, index int) []CardID {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
