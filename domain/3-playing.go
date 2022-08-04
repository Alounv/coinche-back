package domain

import (
	"errors"
	"fmt"
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

func (player Player) hasCard(card cardID) bool {
	for _, cardID := range player.Hand {
		if cardID == card {
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

func (player Player) hasNoBiggerTrump(trump Color, turn turn) bool {
	for _, cardID := range player.Hand {
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

	if player.hasColor(askedColor) {
		return errors.New(ErrShouldPlayAskedColor)
	}

	if player.hasNoTrump(game.trump) {
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

	if lastTurn.isTheBiggestTrump(card, game.trump) || player.hasNoBiggerTrump(game.trump, lastTurn) {
		return nil
	}

	return errors.New(ErrShouldPlayBiggerTrump)
}

func (game *Game) createTurn(newPlay play) {
	game.turns = append(game.turns, turn{
		plays: []play{newPlay},
	})
}

func (game *Game) updateTurn(newPlay play) {
	lastTurnIndex := len(game.turns) - 1
	lastTurn := game.turns[lastTurnIndex]

	lastTurn.plays = append(lastTurn.plays, newPlay)

	if len(lastTurn.plays) == 4 {
		lastTurn.setWinner(game.trump)
		game.setFirstPlayer(lastTurn.winner)
	}

	game.turns[lastTurnIndex] = lastTurn
}

func (game *Game) allCardsPlayed() bool {
	return len(game.turns) == 8 && len(game.turns[7].plays) == 4
}

func (card card) getStrength(trump Color) Strength {
	if trump == card.color || trump == AllTrump {
		return card.TrumpStrength
	} else {
		return card.strength
	}
}

func (game Game) getPlayersCards() map[string][]cardID {
	playersCards := map[string][]cardID{}

	for _, turn := range game.turns {
		for _, play := range turn.plays {
			playersCards[turn.winner] = append(playersCards[turn.winner], play.card)
		}
	}

	return playersCards
}

func getScoreWithCoinche(score int, coinche int) int {
	if coinche > 0 {
		return score * 2 * coinche
	} else {
		return score
	}
}

func (game Game) getTeamPoints() (points map[string]int, scores map[string]int) {
	playersCards := game.getPlayersCards()

	teamPoints := map[string]int{}

	potentialBelotes := map[Color]string{}

	for player, playerCards := range playersCards {
		team := game.Players[player].Team
		potentialPlayerBelotes := map[Color]int{}

		for _, card := range playerCards {
			teamPoints[team] += cards[card].getValue(game.trump)

			card := cards[card]
			cardStrength := card.getStrength(game.trump)
			if cardStrength == TQueen || cardStrength == TKing {
				potentialPlayerBelotes[card.color]++
			}
		}

		for color, number := range potentialPlayerBelotes {
			if number == 2 {
				potentialBelotes[color] = player
			}
		}
	}

	lastTurn := game.turns[len(game.turns)-1]
	lastWinnerTeam := game.Players[lastTurn.winner].Team
	teamPoints[lastWinnerTeam] += 10

	if game.trump == NoTrump {
		for team, points := range teamPoints {
			teamPoints[team] = points * 162 / 130
		}
	} else if game.trump == AllTrump {
		for team, points := range teamPoints {
			teamPoints[team] = points * 162 / 258
		}
	}

	// IF PLAYER TEAM AS TAKEN
	playerWithBelote := ""

	if len(potentialBelotes) > 0 {
		for _, turn := range game.turns {
			for _, play := range turn.plays {
				card := cards[play.card]
				cardStrength := card.getStrength(game.trump)
				if cardStrength == TQueen || cardStrength == TKing {
					if player, ok := potentialBelotes[card.color]; ok {
						playerWithBelote = player
					}
				}
			}
		}
	}

	lastBid, contract := game.getLastBid()
	lastBidTeam := game.Players[lastBid.Player].Team
	otherTeam := ""

	for team := range teamPoints {
		if team != lastBidTeam {
			otherTeam = team
		}
	}

	teamScore := map[string]int{
		lastBidTeam: 0,
		otherTeam:   0,
	}

	if playerWithBelote != "" {
		team := game.Players[playerWithBelote].Team
		teamScore[team] += 20

		if lastBidTeam == team {
			teamPoints[team] += 20
		}
	}

	coinche := lastBid.Coinche

	isCapotWon := contract == Capot && teamPoints[otherTeam] == 0
	isNormalContractWon := contract != Capot && teamPoints[lastBidTeam] >= int(contract)

	if isCapotWon {
		teamScore[lastBidTeam] += getScoreWithCoinche(160, coinche)
	} else if isNormalContractWon {
		teamScore[lastBidTeam] += getScoreWithCoinche(int(contract), coinche)
	} else {
		teamScore[otherTeam] += getScoreWithCoinche(160, coinche)
	}

	return teamPoints, teamScore
}

func (game *Game) end() {
	game.Phase = Counting

	teamsPoints, teamScore := game.getTeamPoints()

	fmt.Println(teamsPoints, teamScore)
}

func (game *Game) Play(playerName string, card cardID) error {
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

	newPlay := play{
		playerName: playerName,
		card:       card,
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
